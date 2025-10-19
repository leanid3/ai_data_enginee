package handlers

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"ai-data-engineer-backend/domain/models"
	"ai-data-engineer-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ! FileService интерфейс для работы с файлами
type FileService interface {
	UploadFile(ctx context.Context, userID, filename string, file io.Reader) (string, error)
	GetFileInfo(ctx context.Context, fileID string) (interface{}, error)
	DeleteFile(ctx context.Context, fileID string) error
	ListFiles(ctx context.Context, userID string, limit, offset int) ([]interface{}, error)
}

// ! FileHandler обработчик для работы с файлами
type FileHandler struct {
	fileService FileService
	logger      logger.Logger
}

// ! NewFileHandler создает новый FileHandler
func NewFileHandler(fileService FileService, logger logger.Logger) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		logger:      logger,
	}
}

// * UploadFile загружает и анализирует файл
func (h *FileHandler) UploadFile(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	requestLogger.Info("Starting: Handler.FileHandler.UploadFile")

	// Получаем файл из multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		requestLogger.WithField("error", err.Error()).Warn("Failed to get file from form")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "file_not_found",
			Message:   "Файл не найден в запросе",
			Timestamp: time.Now(),
		})
		return
	}
	defer file.Close()

	// Получаем тип файла
	fileType := c.PostForm("file_type")
	if fileType != "csv" && fileType != "json" && fileType != "xml" && fileType != "" {
		requestLogger.WithField("error", "invalid_file_type").Warn("Invalid file type")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "invalid_file_type",
			Message:   "Неверный тип файла",
			Timestamp: time.Now(),
		})
		return
	}

	// Получаем ID пользователя
	userID := c.PostForm("user_id")
	if userID == "" {
		userID = "default_user" // По умолчанию
	}

	// Получаем целевую систему
	targetDB := c.PostForm("target_db")
	if targetDB == "" {
		targetDB = "postgres" // По умолчанию
	}

	// Загружаем файл
	result, err := h.fileService.UploadFile(c.Request.Context(), userID, header.Filename, file)
	if err != nil {
		requestLogger.WithField("error", err.Error()).Error("Failed to upload file")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "upload_failed",
			Message:   "Ошибка загрузки файла",
			Timestamp: time.Now(),
		})
		return
	}

	response := models.FileUploadResponse{
		FileID:    result,
		Status:    "uploaded",
		Message:   "Файл успешно загружен и проанализирован",
		CreatedAt: time.Now(),
	}
	requestLogger.WithField("file_id", result).Info("File uploaded successfully")

	c.JSON(http.StatusOK, response)
	requestLogger.Info("End: Handler.FileHandler.UploadFile")
}

// GetFileInfo получает информацию о файле
func (h *FileHandler) GetFileInfo(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	fileID := c.Param("id")

	if fileID == "" {
		requestLogger.Warn("Missing file ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID файла обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	fileInfo, err := h.fileService.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("file_id", fileID).Error("Failed to get file info")
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:     "file_not_found",
			Message:   "Файл не найден",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("file_id", fileID).Info("File info retrieved")
	c.JSON(http.StatusOK, fileInfo)
}

// DeleteFile удаляет файл
func (h *FileHandler) DeleteFile(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	fileID := c.Param("id")

	if fileID == "" {
		requestLogger.Warn("Missing file ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID файла обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	err := h.fileService.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("file_id", fileID).Error("Failed to delete file")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "delete_failed",
			Message:   "Ошибка удаления файла",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("file_id", fileID).Info("File deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Файл успешно удален",
		"file_id": fileID,
	})
}

// ListFiles получает список файлов пользователя
func (h *FileHandler) ListFiles(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	userID := c.Query("user_id")

	if userID == "" {
		requestLogger.Warn("Missing user ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID пользователя обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	// Параметры пагинации
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	files, err := h.fileService.ListFiles(c.Request.Context(), userID, limit, offset)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("user_id", userID).Error("Failed to list files")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "list_failed",
			Message:   "Ошибка получения списка файлов",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("user_id", userID).WithField("count", len(files)).Info("Files listed")
	c.JSON(http.StatusOK, gin.H{
		"files":  files,
		"limit":  limit,
		"offset": offset,
		"count":  len(files),
	})
}
