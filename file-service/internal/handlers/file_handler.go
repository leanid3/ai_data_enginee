package handlers

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"file-service/internal/models"
	"file-service/internal/storage"
	"file-service/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileHandler struct {
	proto.UnimplementedFileServiceServer
	db      *gorm.DB
	storage storage.StorageInterface
}

func NewFileHandler(db *gorm.DB, storage storage.StorageInterface) *FileHandler {
	return &FileHandler{
		db:      db,
		storage: storage,
	}
}

// UploadFile обрабатывает загрузку файла через gRPC
func (h *FileHandler) UploadFile(ctx context.Context, req *proto.UploadFileRequest) (*proto.UploadFileResponse, error) {
	// Генерация уникального ID файла
	fileID := uuid.New().String()

	// Создание записи в базе данных
	file := &models.File{
		ID:          fileID,
		UserID:      req.UserId,
		Filename:    req.Filename,
		ContentType: req.ContentType,
		FileSize:    req.FileSize,
		Status:      "uploaded",
		StoragePath: fmt.Sprintf("users/%s/%s/%s", req.UserId, fileID, req.Filename),
		StorageType: "minio",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.db.Create(file).Error; err != nil {
		return &proto.UploadFileResponse{
			FileId:  fileID,
			Status:  "failed",
			Message: fmt.Sprintf("Ошибка сохранения метаданных: %v", err),
		}, nil
	}

	return &proto.UploadFileResponse{
		FileId:      fileID,
		Status:      "uploaded",
		Message:     "Файл успешно загружен",
		StoragePath: file.StoragePath,
		FileSize:    req.FileSize,
		CreatedAt:   file.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetFileInfo получает информацию о файле
func (h *FileHandler) GetFileInfo(ctx context.Context, req *proto.GetFileInfoRequest) (*proto.GetFileInfoResponse, error) {
	var file models.File
	if err := h.db.Where("id = ? AND user_id = ?", req.FileId, req.UserId).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &proto.GetFileInfoResponse{}, fmt.Errorf("файл не найден")
		}
		return &proto.GetFileInfoResponse{}, err
	}

	// Получение информации из хранилища
	storageInfo, err := h.storage.GetFileInfo(ctx, "files", file.StoragePath)
	if err != nil {
		return &proto.GetFileInfoResponse{}, err
	}

	return &proto.GetFileInfoResponse{
		FileId:      file.ID,
		Filename:    file.Filename,
		ContentType: file.ContentType,
		FileSize:    file.FileSize,
		Status:      file.Status,
		StoragePath: file.StoragePath,
		CreatedAt:   file.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   file.UpdatedAt.Format(time.RFC3339),
		Metadata: map[string]string{
			"storage_type": file.StorageType,
			"checksum":     file.Checksum,
			"etag":         storageInfo.ETag,
		},
	}, nil
}

// DownloadFile скачивает файл
func (h *FileHandler) DownloadFile(req *proto.DownloadFileRequest, stream proto.FileService_DownloadFileServer) error {
	// Проверка существования файла
	var file models.File
	if err := h.db.Where("id = ? AND user_id = ?", req.FileId, req.UserId).First(&file).Error; err != nil {
		return fmt.Errorf("файл не найден")
	}

	// Получение файла из хранилища
	reader, err := h.storage.DownloadFile(stream.Context(), "files", file.StoragePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Отправка файла по частям
	buffer := make([]byte, 64*1024) // 64KB chunks
	offset := int64(0)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		response := &proto.DownloadFileResponse{
			Chunk:         buffer[:n],
			TotalSize:     file.FileSize,
			CurrentOffset: offset,
		}

		if err := stream.Send(response); err != nil {
			return err
		}

		offset += int64(n)
	}

	return nil
}

// DeleteFile удаляет файл
func (h *FileHandler) DeleteFile(ctx context.Context, req *proto.DeleteFileRequest) (*proto.DeleteFileResponse, error) {
	// Проверка существования файла
	var file models.File
	if err := h.db.Where("id = ? AND user_id = ?", req.FileId, req.UserId).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &proto.DeleteFileResponse{
				FileId:  req.FileId,
				Status:  "not_found",
				Message: "Файл не найден",
			}, nil
		}
		return &proto.DeleteFileResponse{}, err
	}

	// Удаление из хранилища
	if err := h.storage.DeleteFile(ctx, "files", file.StoragePath); err != nil {
		return &proto.DeleteFileResponse{
			FileId:  req.FileId,
			Status:  "failed",
			Message: fmt.Sprintf("Ошибка удаления из хранилища: %v", err),
		}, nil
	}

	// Мягкое удаление из базы данных
	if err := h.db.Delete(&file).Error; err != nil {
		return &proto.DeleteFileResponse{
			FileId:  req.FileId,
			Status:  "failed",
			Message: fmt.Sprintf("Ошибка удаления из базы данных: %v", err),
		}, nil
	}

	return &proto.DeleteFileResponse{
		FileId:  req.FileId,
		Status:  "deleted",
		Message: "Файл успешно удален",
	}, nil
}

// ListFiles получает список файлов пользователя
func (h *FileHandler) ListFiles(ctx context.Context, req *proto.ListFilesRequest) (*proto.ListFilesResponse, error) {
	var files []models.File
	query := h.db.Where("user_id = ?", req.UserId)

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Пагинация
	offset := int((req.Page - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).Order("created_at DESC").Find(&files).Error; err != nil {
		return &proto.ListFilesResponse{}, err
	}

	// Подсчет общего количества
	var totalCount int64
	h.db.Model(&models.File{}).Where("user_id = ?", req.UserId).Count(&totalCount)

	// Преобразование в ответ
	var fileInfos []*proto.FileInfo
	for _, file := range files {
		fileInfos = append(fileInfos, &proto.FileInfo{
			FileId:      file.ID,
			Filename:    file.Filename,
			ContentType: file.ContentType,
			FileSize:    file.FileSize,
			Status:      file.Status,
			CreatedAt:   file.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   file.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &proto.ListFilesResponse{
		Files:      fileInfos,
		TotalCount: int32(totalCount),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// GetUploadStatus получает статус загрузки файла
func (h *FileHandler) GetUploadStatus(ctx context.Context, req *proto.GetUploadStatusRequest) (*proto.GetUploadStatusResponse, error) {
	var file models.File
	if err := h.db.Where("id = ? AND user_id = ?", req.FileId, req.UserId).First(&file).Error; err != nil {
		return &proto.GetUploadStatusResponse{}, fmt.Errorf("файл не найден")
	}

	progress := int32(100)
	if file.Status == "uploading" {
		progress = 50 // Примерное значение
	}

	return &proto.GetUploadStatusResponse{
		FileId:          file.ID,
		Status:          file.Status,
		ProgressPercent: progress,
		Message:         "Статус загрузки",
	}, nil
}

// HTTP обработчики для multipart uploads

// RegisterHTTPRoutes регистрирует HTTP маршруты
func (h *FileHandler) RegisterHTTPRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/files/upload", h.handleMultipartUpload)
		api.GET("/files/:file_id", h.handleGetFileInfo)
		api.DELETE("/files/:file_id", h.handleDeleteFile)
		api.GET("/files", h.handleListFiles)
	}
}

// handleMultipartUpload обрабатывает multipart загрузку файлов
func (h *FileHandler) handleMultipartUpload(c *gin.Context) {
	// Получение файла из multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения файла"})
		return
	}
	defer file.Close()

	// Получение user_id из заголовка или параметра
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = c.Query("user_id")
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id обязателен"})
		return
	}

	// Генерация уникального ID файла
	fileID := uuid.New().String()

	// Определение content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Создание записи в базе данных
	dbFile := &models.File{
		ID:          fileID,
		UserID:      userID,
		Filename:    header.Filename,
		ContentType: contentType,
		FileSize:    header.Size,
		Status:      "uploading",
		StoragePath: fmt.Sprintf("users/%s/%s/%s", userID, fileID, header.Filename),
		StorageType: "minio",
		Metadata:    "{}",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Вычисление checksum
	hasher := md5.New()
	file.Seek(0, 0)
	if _, err := io.Copy(hasher, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка вычисления checksum"})
		return
	}
	dbFile.Checksum = fmt.Sprintf("%x", hasher.Sum(nil))

	// Сохранение в базу данных
	if err := h.db.Create(dbFile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения метаданных"})
		return
	}

	// Загрузка в хранилище
	file.Seek(0, 0)
	if err := h.storage.UploadFile(c.Request.Context(), "files", dbFile.StoragePath, file, header.Size, contentType); err != nil {
		// Обновление статуса на failed
		h.db.Model(dbFile).Update("status", "failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки в хранилище"})
		return
	}

	// Обновление статуса на uploaded
	h.db.Model(dbFile).Update("status", "uploaded")

	c.JSON(http.StatusOK, gin.H{
		"file_id":      fileID,
		"status":       "uploaded",
		"message":      "Файл успешно загружен",
		"storage_path": dbFile.StoragePath,
		"file_size":    header.Size,
		"created_at":   dbFile.CreatedAt.Format(time.RFC3339),
	})
}

// handleGetFileInfo получает информацию о файле через HTTP
func (h *FileHandler) handleGetFileInfo(c *gin.Context) {
	fileID := c.Param("file_id")
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = c.Query("user_id")
	}

	req := &proto.GetFileInfoRequest{
		FileId: fileID,
		UserId: userID,
	}

	resp, err := h.GetFileInfo(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// handleDeleteFile удаляет файл через HTTP
func (h *FileHandler) handleDeleteFile(c *gin.Context) {
	fileID := c.Param("file_id")
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = c.Query("user_id")
	}

	req := &proto.DeleteFileRequest{
		FileId: fileID,
		UserId: userID,
	}

	resp, err := h.DeleteFile(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// handleListFiles получает список файлов через HTTP
func (h *FileHandler) handleListFiles(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = c.Query("user_id")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	req := &proto.ListFilesRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Status:   status,
	}

	resp, err := h.ListFiles(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
