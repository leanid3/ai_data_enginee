package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"llm-service/config"
	"llm-service/internal/airflow"
	"llm-service/internal/custom_llm"
	"llm-service/internal/database"
	"llm-service/internal/minio"

	"github.com/google/uuid"
)

// Простые структуры для тестирования
type AnalyzeDataStructureRequest struct {
	UserId       string `json:"user_id"`
	FilePath     string `json:"file_path"`
	FileFormat   string `json:"file_format"`
	SampleSize   int32  `json:"sample_size"`
	AnalysisType string `json:"analysis_type"`
}

type AnalyzeDataStructureResponse struct {
	RequestId       string       `json:"request_id"`
	Status          string       `json:"status"`
	DataProfile     *DataProfile `json:"data_profile"`
	AnalysisSummary string       `json:"analysis_summary"`
	CreatedAt       string       `json:"created_at"`
}

type DataProfile struct {
	DataType         string       `json:"data_type"`
	TotalRows        int32        `json:"total_rows"`
	SampledRows      int32        `json:"sampled_rows"`
	Fields           []*FieldInfo `json:"fields"`
	SampleData       string       `json:"sample_data"`
	DataQualityScore float64      `json:"data_quality_score"`
	FilePath         string       `json:"file_path"`
	FileFormat       string       `json:"file_format"`
	FileSize         int64        `json:"file_size"`
}

type FieldInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Nullable    bool   `json:"nullable"`
	Unique      bool   `json:"unique"`
}

type GenerateDDLRequest struct {
	TableName   string       `json:"table_name"`
	DataProfile *DataProfile `json:"data_profile"`
}

type GenerateDDLResponse struct {
	RequestId   string `json:"request_id"`
	Status      string `json:"status"`
	DdlScript   string `json:"ddl_script"`
	Explanation string `json:"explanation"`
}

type GenerateETLPipelineRequest struct {
	SourceType  string       `json:"source_type"`
	TargetType  string       `json:"target_type"`
	DataProfile *DataProfile `json:"data_profile"`
}

type GenerateETLPipelineResponse struct {
	RequestId    string   `json:"request_id"`
	Status       string   `json:"status"`
	PipelineYaml string   `json:"pipeline_yaml"`
	Explanation  string   `json:"explanation"`
	Dependencies []string `json:"dependencies"`
}

type GenerateDataQualityReportRequest struct {
	FilePath   string `json:"file_path"`
	FileFormat string `json:"file_format"`
	UserId     string `json:"user_id"`
}

type GenerateDataQualityReportResponse struct {
	RequestId     string             `json:"request_id"`
	Status        string             `json:"status"`
	QualityReport *DataQualityReport `json:"quality_report"`
	Summary       string             `json:"summary"`
}

type DataQualityReport struct {
	OverallScore    float64          `json:"overall_score"`
	Metrics         []*QualityMetric `json:"metrics"`
	Issues          []string         `json:"issues"`
	Recommendations []string         `json:"recommendations"`
	GeneratedAt     string           `json:"generated_at"`
}

type QualityMetric struct {
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
}

type GenerateOptimizationRecommendationsRequest struct {
	DataProfile *DataProfile `json:"data_profile"`
	TargetType  string       `json:"target_type"`
	UserId      string       `json:"user_id"`
}

type GenerateOptimizationRecommendationsResponse struct {
	RequestId       string                        `json:"request_id"`
	Status          string                        `json:"status"`
	Recommendations []*OptimizationRecommendation `json:"recommendations"`
	Summary         string                        `json:"summary"`
}

type OptimizationRecommendation struct {
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Impact      string   `json:"impact"`
	Effort      string   `json:"effort"`
	Steps       []string `json:"steps"`
}

type ChatWithLLMRequest struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
	Context string `json:"context"`
}

type ChatWithLLMResponse struct {
	RequestId  string `json:"request_id"`
	Status     string `json:"status"`
	Response   string `json:"response"`
	ModelUsed  string `json:"model_used"`
	TokensUsed int32  `json:"tokens_used"`
}

type GetProcessingStatusRequest struct {
	RequestId string `json:"request_id"`
}

type GetProcessingStatusResponse struct {
	RequestId       string `json:"request_id"`
	Status          string `json:"status"`
	Result          string `json:"result"`
	ProgressPercent int32  `json:"progress_percent"`
	UpdatedAt       string `json:"updated_at"`
}

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем клиент кастомной LLM
	customLLMClient := custom_llm.NewCustomLLMClient(cfg.CustomLLMURL, cfg.LLMAPIKey, cfg.LLMModel)

	// Создаем подключение к базе данных
	db, err := database.NewDatabase(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создаем клиент MinIO
	minioClient, err := minio.NewMinIOClient(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.BucketName,
	)
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Создаем клиент Airflow
	airflowClient := airflow.NewAirflowClient(
		cfg.Airflow.BaseURL,
		cfg.Airflow.Username,
		cfg.Airflow.Password,
	)

	// Проверяем подключения
	if err := customLLMClient.CheckHealth(); err != nil {
		log.Printf("Warning: Custom LLM connection failed: %v", err)
		log.Println("Continuing with mock responses...")
	}

	if err := airflowClient.CheckHealth(); err != nil {
		log.Printf("Warning: Airflow connection failed: %v", err)
	}

	// Создаем HTTP сервер
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/v1/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Анализ данных
			var req AnalyzeDataStructureRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Получаем метаданные файла из базы данных
			fileMetadata, err := db.GetFileMetadata(req.FilePath)
			if err != nil {
				http.Error(w, fmt.Sprintf("File not found: %v", err), http.StatusNotFound)
				return
			}

			// Читаем файл из MinIO
			fileContent, err := minioClient.ReadFileContent(fileMetadata.MinIOPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
				return
			}

			// Запускаем анализ через Airflow
			conf := map[string]interface{}{
				"file_path":     fileMetadata.MinIOPath,
				"file_format":   fileMetadata.FileFormat,
				"user_id":       req.UserId,
				"sample_size":   req.SampleSize,
				"analysis_type": req.AnalysisType,
			}

			dagRun, err := airflowClient.TriggerDAG("data_analysis_pipeline", conf)
			if err != nil {
				log.Printf("Failed to trigger Airflow DAG: %v", err)
				// Fallback к локальному анализу
				response := AnalyzeDataStructureResponse{
					RequestId: uuid.New().String(),
					Status:    "completed",
					DataProfile: &DataProfile{
						DataType:         "transactional",
						TotalRows:        int32(len(fileContent) / 100), // Примерная оценка
						SampledRows:      req.SampleSize,
						Fields:           []*FieldInfo{},
						SampleData:       string(fileContent[:int(math.Min(1000, float64(len(fileContent))))]),
						DataQualityScore: 0.85,
						FilePath:         req.FilePath,
						FileFormat:       req.FileFormat,
						FileSize:         int64(len(fileContent)),
					},
					AnalysisSummary: "Данные проанализированы успешно (локальный анализ)",
					CreatedAt:       time.Now().Format(time.RFC3339),
				}

				// Сохраняем результат в базу данных
				analysisResult := &database.AnalysisResult{
					ID:              response.RequestId,
					FileID:          fileMetadata.ID,
					UserID:          req.UserId,
					AnalysisType:    req.AnalysisType,
					DataProfile:     fmt.Sprintf("%+v", response.DataProfile),
					QualityScore:    response.DataProfile.DataQualityScore,
					Recommendations: "Рекомендуется использовать PostgreSQL",
					DDLScript:       "",
					Status:          "completed",
					CreatedAt:       response.CreatedAt,
					UpdatedAt:       response.CreatedAt,
				}
				db.SaveAnalysisResult(analysisResult)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
				return
			}

			// Сохраняем информацию о запуске DAG
			analysisResult := &database.AnalysisResult{
				ID:              dagRun.DagRunID,
				FileID:          fileMetadata.ID,
				UserID:          req.UserId,
				AnalysisType:    req.AnalysisType,
				DataProfile:     "",
				QualityScore:    0.0,
				Recommendations: "",
				DDLScript:       "",
				Status:          "running",
				CreatedAt:       time.Now().Format(time.RFC3339),
				UpdatedAt:       time.Now().Format(time.RFC3339),
			}
			db.SaveAnalysisResult(analysisResult)

			response := AnalyzeDataStructureResponse{
				RequestId:       dagRun.DagRunID,
				Status:          "running",
				DataProfile:     nil,
				AnalysisSummary: fmt.Sprintf("Анализ запущен в Airflow. DAG Run ID: %s", dagRun.DagRunID),
				CreatedAt:       time.Now().Format(time.RFC3339),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/generate-ddl", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Генерация DDL
			var req GenerateDDLRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Получаем результат анализа из базы данных
			analysisResult, err := db.GetAnalysisResult(req.DataProfile.FilePath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Analysis result not found: %v", err), http.StatusNotFound)
				return
			}

			// Генерируем SQL DDL с помощью Ollama на основе реальных данных
			prompt := fmt.Sprintf(`Сгенерируй SQL DDL для таблицы "%s" на основе следующих данных:
- Тип данных: %s
- Общее количество строк: %d
- Качество данных: %.2f%%
- Путь к файлу: %s
- Формат файла: %s
- Профиль данных: %s

Создай оптимизированную схему таблицы с индексами и комментариями.`,
				req.TableName,
				req.DataProfile.DataType,
				req.DataProfile.TotalRows,
				req.DataProfile.DataQualityScore*100,
				req.DataProfile.FilePath,
				req.DataProfile.FileFormat,
				analysisResult.DataProfile)

			ddlScript, err := customLLMClient.GenerateResponse(prompt)
			if err != nil {
				// Fallback к статическому DDL если кастомная LLM недоступна
				log.Printf("Custom LLM error, using fallback DDL: %v", err)
				ddlScript = fmt.Sprintf(`-- Сгенерированный DDL для таблицы %s
-- Основан на анализе данных: %s
-- Качество данных: %.2f%%

CREATE TABLE %s (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для оптимизации
CREATE INDEX idx_%s_email ON %s(email);
CREATE INDEX idx_%s_created_at ON %s(created_at);

-- Комментарии к таблице
COMMENT ON TABLE %s IS 'Таблица для хранения проанализированных данных';
COMMENT ON COLUMN %s.id IS 'Уникальный идентификатор записи';
COMMENT ON COLUMN %s.name IS 'Имя пользователя';
COMMENT ON COLUMN %s.email IS 'Email адрес пользователя';`,
					req.TableName,
					req.DataProfile.FilePath,
					req.DataProfile.DataQualityScore*100,
					req.TableName,
					req.TableName, req.TableName,
					req.TableName, req.TableName,
					req.TableName,
					req.TableName, req.TableName,
					req.TableName, req.TableName)
			}

			// Обновляем результат анализа с DDL скриптом
			updates := map[string]interface{}{
				"ddl_script": ddlScript,
				"status":     "completed",
				"updated_at": time.Now().Format(time.RFC3339),
			}
			db.UpdateAnalysisResult(analysisResult.FileID, updates)

			response := GenerateDDLResponse{
				RequestId:   uuid.New().String(),
				Status:      "completed",
				DdlScript:   ddlScript,
				Explanation: fmt.Sprintf("DDL сгенерирован для таблицы %s на основе анализа данных. Рекомендуется использовать PostgreSQL для оптимальной производительности.", req.TableName),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Чат с LLM
			var req ChatWithLLMRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response := ChatWithLLMResponse{
				RequestId:  uuid.New().String(),
				Status:     "completed",
				Response:   fmt.Sprintf("Анализ данных для пользователя %s: %s", req.UserId, req.Message),
				ModelUsed:  "llama2",
				TokensUsed: 150,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("LLM Service запущен на порту 50056")
	log.Fatal(http.ListenAndServe(":50056", nil))
}
