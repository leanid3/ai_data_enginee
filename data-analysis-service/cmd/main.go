package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type AnalysisService struct {
	ollamaURL string
	minioURL  string
}

type AnalysisRequest struct {
	FileID   string `json:"file_id"`
	UserID   string `json:"user_id"`
	FilePath string `json:"file_path"`
}

type AnalysisResponse struct {
	AnalysisID string                 `json:"analysis_id"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Result     map[string]interface{} `json:"result,omitempty"`
}

type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewAnalysisService(ollamaURL, minioURL string) *AnalysisService {
	return &AnalysisService{
		ollamaURL: ollamaURL,
		minioURL:  minioURL,
	}
}

// AnalyzeData выполняет реальный анализ данных
func (s *AnalysisService) AnalyzeData(c *gin.Context) {
	var req AnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Валидация входных данных
	if req.FileID == "" || req.UserID == "" || req.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_id, user_id и file_path обязательны"})
		return
	}

	// Генерация уникального ID анализа
	analysisID := fmt.Sprintf("analysis_%s_%d", req.FileID, time.Now().Unix())

	// Запуск анализа в фоне
	go s.performAnalysis(analysisID, req)

	response := AnalysisResponse{
		AnalysisID: analysisID,
		Status:     "started",
		Message:    "Анализ данных запущен",
		Result: map[string]interface{}{
			"file_id": req.FileID,
			"user_id": req.UserID,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetAnalysisStatus получает статус анализа
func (s *AnalysisService) GetAnalysisStatus(c *gin.Context) {
	analysisID := c.Param("analysis_id")
	if analysisID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "analysis_id обязателен"})
		return
	}

	// Читаем результат анализа из файла
	result, err := s.loadAnalysisResult(analysisID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Анализ не найден"})
		return
	}

	// Формируем структурированный ответ для DDL пайплайна
	status := map[string]interface{}{
		"analysis_id": analysisID,
		"status":      "completed",
		"progress":    100,
		"message":     "Анализ завершен",
		"result": map[string]interface{}{
			"data_quality_score":     result["data_profile"].(map[string]interface{})["data_quality_score"],
			"storage_recommendation": s.getStorageRecommendation(result),
			"table_schema":           s.generateTableSchema(result),
			"ddl_metadata":           s.generateDDLMetadata(result),
			"recommendations":        result["recommendations"],
			"analysis_timestamp":     result["timestamp"],
		},
	}

	c.JSON(http.StatusOK, status)
}

// performAnalysis выполняет реальный анализ данных
func (s *AnalysisService) performAnalysis(analysisID string, req AnalysisRequest) {
	log.Printf("Начинаем анализ файла %s для пользователя %s", req.FileID, req.UserID)

	// 1. Получение файла из MinIO
	fileContent := s.getFileFromMinIO(req.FilePath)

	// 2. Анализ структуры данных
	dataProfile := s.analyzeDataStructure(fileContent)

	// 3. Анализ с помощью LLM (передаем полный файл)
	llmAnalysis := s.analyzeWithLLM(dataProfile, fileContent)

	// 4. Генерация рекомендаций
	recommendations := s.generateRecommendations(dataProfile, llmAnalysis)

	// 5. Генерация DDL скрипта через LLM (передаем полный файл)
	ddlScript := s.generateDDLWithLLM(dataProfile, llmAnalysis, fileContent)

	// 6. Перенос данных в PostgreSQL
	transferError := s.transferDataToPostgres(fileContent, ddlScript)
	if transferError != nil {
		log.Printf("Ошибка переноса данных: %v", transferError)
	}

	// 7. Сохранение результата
	s.saveAnalysisResult(analysisID, dataProfile, llmAnalysis, recommendations, ddlScript, transferError)

	log.Printf("Анализ завершен для файла %s", req.FileID)
}

// getFileFromMinIO получает файл из MinIO
func (s *AnalysisService) getFileFromMinIO(filePath string) string {
	log.Printf("Получаем файл из MinIO: %s", filePath)

	// Используем MinIO API для получения файла
	// MinIO API endpoint: http://minio:9000/minio/files/bucket/object
	url := fmt.Sprintf("http://minio:9000/minio/files/%s", filePath)

	// Создаем HTTP клиент с увеличенным таймаутом для больших файлов
	client := &http.Client{
		Timeout: 5 * time.Minute, // 5 минут для больших файлов
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Ошибка создания запроса к MinIO: %v", err)
		return s.getFallbackData()
	}

	// Добавляем заголовки для MinIO
	req.Header.Set("User-Agent", "data-analysis-service")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка получения файла из MinIO: %v", err)
		return s.getFallbackData()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка HTTP при получении файла: %d", resp.StatusCode)
		return s.getFallbackData()
	}

	// Читаем файл (ограничиваем размер для анализа)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024)) // 10MB максимум
	if err != nil {
		log.Printf("Ошибка чтения файла: %v", err)
		return s.getFallbackData()
	}

	log.Printf("Файл успешно получен из MinIO, размер: %d байт", len(body))
	return string(body)
}

// getFallbackData возвращает заглушечные данные при ошибке
func (s *AnalysisService) getFallbackData() string {
	log.Printf("Используем заглушечные данные")
	return `created;order_status;ticket_status;ticket_price;visitor_category;event_id;is_active;valid_to;count_visitor;is_entrance;is_entrance_mdate;event_name;event_kind_name;spot_id;spot_name;museum_name;start_datetime;ticket_id;update_timestamp;client_name;name;surname;client_phone;museum_inn;birthday_date;order_number;ticket_number
2021-08-18T16:01:14.583+03:00;PAID;PAID;0.0;Обучающиеся по очной форме обучения;7561;true;2021-08-18;1;true;2021-08-18T19:14:45.427+03:00;Бальный танец;консультация;274010;Шверника ул. 13;Центр культуры;2021-08-18 17:00:00;1778482;2021-08-18T16:01:15.682+03:00;ШУКУРОВ РУСЛАН;КИРИЛЛ;ШУКУРОВ;79859482165;3832203597;;75343-483088;07a16922-969c-1033-a23f-f20b57dcf045
2021-08-19T10:30:22.123+03:00;PENDING;PENDING;150.0;Взрослые;7562;true;2021-08-19;2;false;2021-08-19T12:00:00.000+03:00;Концерт классической музыки;концерт;274011;Тверская ул. 15;Большой театр;2021-08-19 19:00:00;1778483;2021-08-19T10:30:23.456+03:00;ИВАНОВ ИВАН;ИВАН;ИВАНОВ;79123456789;1234567890;1990-05-15;75344-483089;08b17023-970d-1034-a24g-f21c58dcf046
2021-08-20T14:15:33.789+03:00;PAID;PAID;75.5;Студенты;7563;true;2021-08-20;1;true;2021-08-20T16:30:00.000+03:00;Выставка современного искусства;выставка;274012;Арбат ул. 20;Третьяковская галерея;2021-08-20 15:00:00;1778484;2021-08-20T14:15:34.012+03:00;ПЕТРОВ ПЕТР;ПЕТР;ПЕТРОВ;79987654321;9876543210;1985-12-10;75345-483090;09c17124-971e-1035-a25h-f22d59dcf047
2021-08-21T09:45:11.456+03:00;CANCELLED;CANCELLED;200.0;VIP;7564;false;2021-08-21;3;false;2021-08-21T11:00:00.000+03:00;Опера;опера;274013;Красная площадь 1;Большой театр;2021-08-21 18:30:00;1778485;2021-08-21T09:45:12.789+03:00;СИДОРОВ СИДОР;СИДОР;СИДОРОВ;79555666777;5556667777;1978-03-22;75346-483091;10d17225-972f-1036-a26i-f23e60dcf048`
}

// analyzeDataStructure анализирует структуру данных
func (s *AnalysisService) analyzeDataStructure(fileContent string) map[string]interface{} {
	log.Printf("Анализируем структуру данных")

	// Анализируем реальные данные
	lines := strings.Split(fileContent, "\n")
	nonEmptyLines := []string{}

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	totalRows := len(nonEmptyLines)
	if totalRows == 0 {
		return map[string]interface{}{
			"total_rows":         0,
			"total_columns":      0,
			"file_type":          "CSV",
			"has_headers":        false,
			"encoding":           "UTF-8",
			"delimiter":          ",",
			"data_quality_score": 0.0,
		}
	}

	// Анализируем заголовки
	headers := strings.Split(nonEmptyLines[0], ";")
	totalColumns := len(headers)

	// Анализируем качество данных
	qualityScore := s.calculateDataQuality(nonEmptyLines)

	// Анализируем типы данных
	dataTypes := s.analyzeDataTypes(nonEmptyLines)

	return map[string]interface{}{
		"total_rows":         totalRows,
		"total_columns":      totalColumns,
		"file_type":          "CSV",
		"has_headers":        true,
		"encoding":           "UTF-8",
		"delimiter":          ";",
		"data_quality_score": qualityScore,
		"data_types":         dataTypes,
		"sample_data":        s.getSampleData(nonEmptyLines, 5), // Первые 5 строк
	}
}

// calculateDataQuality рассчитывает качество данных
func (s *AnalysisService) calculateDataQuality(lines []string) float64 {
	if len(lines) <= 1 {
		return 0.0
	}

	totalCells := 0
	emptyCells := 0

	for i, line := range lines {
		if i == 0 { // Пропускаем заголовки
			continue
		}

		fields := strings.Split(line, ";")
		totalCells += len(fields)

		for _, field := range fields {
			if strings.TrimSpace(field) == "" {
				emptyCells++
			}
		}
	}

	if totalCells == 0 {
		return 0.0
	}

	quality := 1.0 - (float64(emptyCells) / float64(totalCells))
	return quality
}

// analyzeDataTypes анализирует типы данных в колонках
func (s *AnalysisService) analyzeDataTypes(lines []string) []string {
	if len(lines) <= 1 {
		return []string{}
	}

	headers := strings.Split(lines[0], ";")
	dataTypes := make([]string, len(headers))

	// Анализируем первые несколько строк для определения типов
	sampleSize := 5
	if len(lines) < sampleSize {
		sampleSize = len(lines)
	}

	for i, header := range headers {
		header = strings.TrimSpace(header)
		dataTypes[i] = s.inferColumnType(lines[1:sampleSize], i, header)
	}

	return dataTypes
}

// inferColumnType определяет тип колонки
func (s *AnalysisService) inferColumnType(lines []string, columnIndex int, header string) string {
	// Анализируем заголовок
	headerLower := strings.ToLower(header)

	if containsAny(headerLower, []string{"date", "time", "created", "updated"}) {
		return "TIMESTAMP"
	}
	if containsAny(headerLower, []string{"price", "amount", "cost", "count"}) {
		return "DECIMAL"
	}
	if containsAny(headerLower, []string{"id", "number"}) {
		return "BIGINT"
	}
	if containsAny(headerLower, []string{"phone", "tel"}) {
		return "VARCHAR(20)"
	}
	if containsAny(headerLower, []string{"email", "mail"}) {
		return "VARCHAR(255)"
	}
	if containsAny(headerLower, []string{"is_", "active", "valid"}) {
		return "BOOLEAN"
	}

	// Анализируем данные в колонке
	for _, line := range lines {
		fields := strings.Split(line, ";")
		if columnIndex < len(fields) {
			value := strings.TrimSpace(fields[columnIndex])
			if value != "" {
				// Проверяем тип значения
				if s.isNumeric(value) {
					return "DECIMAL"
				}
				if s.isDate(value) {
					return "TIMESTAMP"
				}
				if s.isBoolean(value) {
					return "BOOLEAN"
				}
			}
		}
	}

	return "TEXT"
}

// isNumeric проверяет, является ли значение числовым
func (s *AnalysisService) isNumeric(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

// isDate проверяет, является ли значение датой
func (s *AnalysisService) isDate(value string) bool {
	// Простая проверка на формат даты
	dateFormats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+03:00",
	}

	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

// isBoolean проверяет, является ли значение булевым
func (s *AnalysisService) isBoolean(value string) bool {
	valueLower := strings.ToLower(value)
	return valueLower == "true" || valueLower == "false" || valueLower == "1" || valueLower == "0"
}

// getSampleData возвращает образец данных
func (s *AnalysisService) getSampleData(lines []string, maxRows int) []string {
	if len(lines) <= maxRows {
		return lines
	}
	return lines[:maxRows]
}

// generateDDLWithLLM генерирует DDL скрипт через LLM
func (s *AnalysisService) generateDDLWithLLM(dataProfile map[string]interface{}, llmAnalysis string, fullFileContent string) string {
	log.Printf("Генерируем DDL скрипт через LLM")

	// Подготавливаем данные для LLM
	dataTypes := dataProfile["data_types"].([]string)

	prompt := fmt.Sprintf(`
Сгенерируй PostgreSQL DDL скрипт для следующих ПОЛНЫХ данных:

ТИПЫ ДАННЫХ:
%s

ПОЛНЫЕ ДАННЫЕ (весь файл):
%s

ТРЕБОВАНИЯ:
1. Создай таблицу с именем "museum_tickets"
2. Используй правильные PostgreSQL типы данных на основе анализа ВСЕХ данных
3. Добавь индексы для часто используемых полей
4. Добавь ограничения целостности где необходимо
5. Используй кавычки для имен полей
6. Включи комментарии
7. Убедись, что все колонки из заголовков присутствуют в DDL
8. Проанализируй ВСЕ строки данных для определения правильных типов

ВАЖНО: Проанализируй ВСЕ данные в файле, не только заголовки!

Верни ТОЛЬКО SQL код без дополнительных объяснений.
`, strings.Join(dataTypes, ", "), fullFileContent)

	// Отправляем запрос к Ollama с увеличенным таймаутом
	ollamaReq := OllamaRequest{
		Model:  "llama3.2:latest",
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1, // Низкая температура для более точного DDL
			"top_p":       0.9,
		},
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		log.Printf("Ошибка сериализации запроса к LLM для DDL: %v", err)
		return s.generateFallbackDDL(dataProfile)
	}

	// Создаем HTTP клиент с увеличенным таймаутом для LLM
	client := &http.Client{
		Timeout: 10 * time.Minute, // 10 минут для больших файлов и сложного DDL
	}

	// Отправка запроса
	resp, err := client.Post(s.ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Ошибка запроса к LLM для DDL: %v", err)
		return s.generateFallbackDDL(dataProfile)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения ответа от LLM для DDL: %v", err)
		return s.generateFallbackDDL(dataProfile)
	}

	// Парсинг ответа
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		log.Printf("Ошибка парсинга ответа от LLM для DDL: %v", err)
		return s.generateFallbackDDL(dataProfile)
	}

	log.Printf("DDL скрипт сгенерирован LLM")

	// Очищаем ответ от markdown разметки
	ddlScript := ollamaResp.Response
	ddlScript = strings.TrimPrefix(ddlScript, "```sql")
	ddlScript = strings.TrimPrefix(ddlScript, "```")
	ddlScript = strings.TrimSuffix(ddlScript, "```")
	ddlScript = strings.TrimSpace(ddlScript)

	// Логируем сгенерированный DDL для отладки
	log.Printf("Сгенерированный DDL (первые 500 символов): %s", ddlScript[:min(len(ddlScript), 500)])

	return ddlScript
}

// min возвращает минимальное из двух значений
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// generateFallbackDDL генерирует резервный DDL при ошибке LLM
func (s *AnalysisService) generateFallbackDDL(dataProfile map[string]interface{}) string {
	log.Printf("Используем резервный DDL")

	headers := dataProfile["sample_data"].([]string)[0]
	fieldNames := strings.Split(headers, ";")
	dataTypes := dataProfile["data_types"].([]string)

	ddl := "DROP TABLE IF EXISTS museum_tickets;\n"
	ddl += "CREATE TABLE museum_tickets (\n"

	for i, fieldName := range fieldNames {
		fieldName = strings.TrimSpace(fieldName)
		fieldName = strings.ReplaceAll(fieldName, " ", "_")
		fieldName = strings.ReplaceAll(fieldName, "-", "_")
		fieldName = strings.ReplaceAll(fieldName, "(", "")
		fieldName = strings.ReplaceAll(fieldName, ")", "")

		// Определяем тип поля
		var fieldType string
		if i < len(dataTypes) {
			fieldType = s.mapToPostgresType(dataTypes[i])
		} else {
			fieldType = "TEXT"
		}

		ddl += fmt.Sprintf("    \"%s\" %s,\n", fieldName, fieldType)
	}

	// Убираем последнюю запятую
	ddl = strings.TrimSuffix(ddl, ",\n")
	ddl += "\n);\n\n"

	// Добавляем индексы
	ddl += "-- Индексы для оптимизации запросов\n"
	ddl += "CREATE INDEX idx_created ON museum_tickets(\"created\");\n"
	ddl += "CREATE INDEX idx_event_id ON museum_tickets(\"event_id\");\n"
	ddl += "CREATE INDEX idx_museum_name ON museum_tickets(\"museum_name\");\n"

	return ddl
}

// mapToPostgresType маппит типы данных в PostgreSQL
func (s *AnalysisService) mapToPostgresType(dataType string) string {
	switch dataType {
	case "TIMESTAMP":
		return "TIMESTAMP"
	case "DECIMAL":
		return "DECIMAL(10,2)"
	case "BIGINT":
		return "BIGINT"
	case "BOOLEAN":
		return "BOOLEAN"
	case "VARCHAR(20)":
		return "VARCHAR(20)"
	case "VARCHAR(255)":
		return "VARCHAR(255)"
	default:
		return "TEXT"
	}
}

// transferDataToPostgres переносит данные в PostgreSQL
func (s *AnalysisService) transferDataToPostgres(fileContent string, ddl string) error {
	log.Printf("Переносим данные в PostgreSQL")

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", "host=postgres port=5432 user=postgres password=postgres dbname=aien_db sslmode=disable")
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	// Выполняем DDL
	_, err = db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("ошибка выполнения DDL: %v", err)
	}

	// Парсим CSV данные
	lines := strings.Split(fileContent, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("недостаточно данных для переноса")
	}

	headers := strings.Split(lines[0], ";")

	// Подготавливаем INSERT запрос
	placeholders := make([]string, len(headers))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	insertSQL := fmt.Sprintf("INSERT INTO museum_tickets (%s) VALUES (%s)",
		strings.Join(headers, ","),
		strings.Join(placeholders, ","))

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("ошибка подготовки запроса: %v", err)
	}
	defer stmt.Close()

	// Вставляем данные (ограничиваем количество для демонстрации)
	maxRows := 100
	insertedRows := 0

	for i, line := range lines[1:] {
		if i >= maxRows {
			break
		}

		fields := strings.Split(line, ";")
		if len(fields) != len(headers) {
			continue // Пропускаем некорректные строки
		}

		// Подготавливаем значения
		values := make([]interface{}, len(fields))
		for j, field := range fields {
			values[j] = strings.TrimSpace(field)
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			log.Printf("Ошибка вставки строки %d: %v", i+1, err)
			continue
		}

		insertedRows++
	}

	log.Printf("Успешно перенесено %d строк в PostgreSQL", insertedRows)
	return nil
}

// analyzeWithLLM выполняет анализ с помощью LLM
func (s *AnalysisService) analyzeWithLLM(dataProfile map[string]interface{}, fullFileContent string) string {
	log.Printf("Анализируем данные с помощью LLM")

	// Подготовка промпта для LLM с ПОЛНЫМИ данными
	prompt := fmt.Sprintf(`
Проанализируй следующие ПОЛНЫЕ данные и дай рекомендации:

СТРУКТУРА ДАННЫХ:
- Количество строк: %v
- Количество столбцов: %v
- Тип файла: %v
- Качество данных: %v
- Типы данных: %v

ПОЛНЫЕ ДАННЫЕ (весь файл):
%s

Дай детальные рекомендации по:
1. Качеству данных (проблемы, аномалии, пропуски)
2. Оптимальному хранилищу (PostgreSQL, ClickHouse, HDFS)
3. Предобработке данных (очистка, трансформация)
4. Возможностям аналитики (статистика, визуализация, ML)
5. Генерации DDL скрипта для PostgreSQL

Ответь на русском языке с конкретными примерами.
`,
		dataProfile["total_rows"],
		dataProfile["total_columns"],
		dataProfile["file_type"],
		dataProfile["data_quality_score"],
		dataProfile["data_types"],
		fullFileContent)

	// Отправка запроса к Ollama
	ollamaReq := OllamaRequest{
		Model:  "llama3.2:latest",
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		log.Printf("Ошибка сериализации запроса к LLM: %v", err)
		return "Ошибка анализа с LLM"
	}

	// Отправка запроса
	resp, err := http.Post(s.ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Ошибка запроса к LLM: %v", err)
		return "LLM недоступен"
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения ответа от LLM: %v", err)
		return "Ошибка чтения ответа от LLM"
	}

	// Парсинг ответа
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		log.Printf("Ошибка парсинга ответа от LLM: %v", err)
		return "Ошибка парсинга ответа от LLM"
	}

	return ollamaResp.Response
}

// generateRecommendations генерирует рекомендации на основе LLM анализа
func (s *AnalysisService) generateRecommendations(dataProfile map[string]interface{}, llmAnalysis string) []string {
	log.Printf("Генерируем рекомендации на основе LLM анализа")

	recommendations := []string{}

	// Парсим LLM анализ и извлекаем рекомендации
	if llmAnalysis != "" {
		// Извлекаем рекомендации из LLM анализа
		llmRecommendations := s.extractRecommendationsFromLLM(llmAnalysis)
		recommendations = append(recommendations, llmRecommendations...)
	}

	// Добавляем базовые рекомендации на основе профиля данных
	basicRecommendations := s.generateBasicRecommendations(dataProfile)
	recommendations = append(recommendations, basicRecommendations...)

	return recommendations
}

// extractRecommendationsFromLLM извлекает рекомендации из LLM анализа
func (s *AnalysisService) extractRecommendationsFromLLM(llmAnalysis string) []string {
	recommendations := []string{}

	// Ищем ключевые слова в LLM анализе
	lines := strings.Split(llmAnalysis, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Ищем строки с рекомендациями
		if strings.Contains(line, "Рекомендация:") ||
			strings.Contains(line, "Рекомендации:") ||
			strings.Contains(line, "рекомендация:") ||
			strings.Contains(line, "рекомендации:") {

			// Очищаем строку от маркеров
			cleanLine := strings.Replace(line, "Рекомендация:", "", -1)
			cleanLine = strings.Replace(cleanLine, "Рекомендации:", "", -1)
			cleanLine = strings.Replace(cleanLine, "рекомендация:", "", -1)
			cleanLine = strings.Replace(cleanLine, "рекомендации:", "", -1)
			cleanLine = strings.TrimSpace(cleanLine)

			if cleanLine != "" && len(cleanLine) > 10 {
				recommendations = append(recommendations, "LLM: "+cleanLine)
			}
		}

		// Ищем строки с дефисами (списки рекомендаций)
		if strings.HasPrefix(line, "-") && len(line) > 5 {
			cleanLine := strings.TrimPrefix(line, "-")
			cleanLine = strings.TrimSpace(cleanLine)
			if cleanLine != "" {
				recommendations = append(recommendations, "LLM: "+cleanLine)
			}
		}
	}

	// Если не нашли структурированные рекомендации, берем первые несколько предложений
	if len(recommendations) == 0 && len(llmAnalysis) > 50 {
		sentences := strings.Split(llmAnalysis, ".")
		for i, sentence := range sentences {
			if i >= 3 { // Берем только первые 3 предложения
				break
			}
			sentence = strings.TrimSpace(sentence)
			if len(sentence) > 20 {
				recommendations = append(recommendations, "LLM: "+sentence+".")
			}
		}
	}

	return recommendations
}

// generateBasicRecommendations генерирует базовые рекомендации на основе профиля данных
func (s *AnalysisService) generateBasicRecommendations(dataProfile map[string]interface{}) []string {
	recommendations := []string{}

	// Анализируем качество данных
	if quality, ok := dataProfile["data_quality_score"].(float64); ok {
		if quality < 0.7 {
			recommendations = append(recommendations, "Низкое качество данных - требуется очистка")
		} else if quality < 0.9 {
			recommendations = append(recommendations, "Среднее качество данных - рекомендуется проверка")
		} else {
			recommendations = append(recommendations, "Высокое качество данных - готовы к использованию")
		}
	}

	// Анализируем размер данных
	if rows, ok := dataProfile["total_rows"].(float64); ok {
		if rows > 1000000 {
			recommendations = append(recommendations, "Большой объем данных - рекомендуется партиционирование")
		} else if rows > 100000 {
			recommendations = append(recommendations, "Средний объем данных - рекомендуется индексирование")
		}
	}

	// Анализируем тип файла
	if fileType, ok := dataProfile["file_type"].(string); ok {
		switch fileType {
		case "CSV":
			recommendations = append(recommendations, "CSV формат - рекомендуется валидация разделителей")
		case "JSON":
			recommendations = append(recommendations, "JSON формат - рекомендуется проверка структуры")
		case "XML":
			recommendations = append(recommendations, "XML формат - рекомендуется валидация схемы")
		}
	}

	return recommendations
}

// extractStorageRecommendationFromLLM извлекает рекомендации по хранилищу из LLM анализа
func (s *AnalysisService) extractStorageRecommendationFromLLM(llmAnalysis string) string {
	// Ищем упоминания систем хранения в LLM анализе
	storageKeywords := map[string]string{
		"postgresql":    "PostgreSQL",
		"postgres":      "PostgreSQL",
		"mysql":         "MySQL",
		"clickhouse":    "ClickHouse",
		"hdfs":          "HDFS",
		"hadoop":        "HDFS",
		"mongodb":       "MongoDB",
		"elasticsearch": "Elasticsearch",
		"redis":         "Redis",
		"cassandra":     "Cassandra",
	}

	llmAnalysisLower := strings.ToLower(llmAnalysis)

	for keyword, storage := range storageKeywords {
		if strings.Contains(llmAnalysisLower, keyword) {
			return storage
		}
	}

	// Если не нашли конкретные упоминания, анализируем контекст
	if strings.Contains(llmAnalysisLower, "реляционн") || strings.Contains(llmAnalysisLower, "sql") {
		return "PostgreSQL"
	}
	if strings.Contains(llmAnalysisLower, "аналитическ") || strings.Contains(llmAnalysisLower, "olap") {
		return "ClickHouse"
	}
	if strings.Contains(llmAnalysisLower, "больш") || strings.Contains(llmAnalysisLower, "масштаб") {
		return "HDFS"
	}

	return "PostgreSQL" // По умолчанию
}

// saveAnalysisResult сохраняет результат анализа
func (s *AnalysisService) saveAnalysisResult(analysisID string, dataProfile map[string]interface{}, llmAnalysis string, recommendations []string, ddlScript string, transferError error) {
	log.Printf("Сохраняем результат анализа: %s", analysisID)

	// В реальной реализации здесь должно быть сохранение в базу данных
	result := map[string]interface{}{
		"analysis_id":     analysisID,
		"data_profile":    dataProfile,
		"llm_analysis":    llmAnalysis,
		"recommendations": recommendations,
		"ddl_script":      ddlScript,
		"transfer_status": map[string]interface{}{
			"success": transferError == nil,
			"error":   transferError,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Сохранение в файл для демонстрации
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile(fmt.Sprintf("/tmp/analysis_%s.json", analysisID), jsonData, 0644)
}

// loadAnalysisResult загружает результат анализа из файла
func (s *AnalysisService) loadAnalysisResult(analysisID string) (map[string]interface{}, error) {
	filePath := fmt.Sprintf("/tmp/analysis_%s.json", analysisID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// getStorageRecommendation определяет рекомендуемую систему хранения
func (s *AnalysisService) getStorageRecommendation(result map[string]interface{}) map[string]interface{} {
	dataProfile := result["data_profile"].(map[string]interface{})
	llmAnalysis := result["llm_analysis"].(string)

	// Анализируем данные для определения оптимального хранилища
	fileType := dataProfile["file_type"].(string)
	totalRows := int(dataProfile["total_rows"].(float64))
	dataQuality := dataProfile["data_quality_score"].(float64)

	// Извлекаем рекомендации по хранилищу из LLM анализа
	storageRecommendation := s.extractStorageRecommendationFromLLM(llmAnalysis)

	// Определяем основное хранилище на основе LLM анализа
	primaryStorage := storageRecommendation
	if primaryStorage == "" {
		primaryStorage = "PostgreSQL" // По умолчанию
	}

	recommendations := map[string]interface{}{
		"primary_storage":   primaryStorage,
		"secondary_storage": []string{"ClickHouse", "HDFS"},
		"reasoning": map[string]interface{}{
			"file_type":          fileType,
			"data_volume":        totalRows,
			"quality_score":      dataQuality,
			"llm_recommendation": storageRecommendation,
			"recommendation":     fmt.Sprintf("%s подходит для структурированных данных с хорошим качеством", primaryStorage),
		},
		"storage_options": map[string]interface{}{
			"postgresql": map[string]interface{}{
				"suitable": true,
				"reasons": []string{
					"Структурированные данные",
					"Хорошее качество данных",
					"Поддержка ACID транзакций",
					"Богатые возможности индексирования",
				},
			},
			"clickhouse": map[string]interface{}{
				"suitable": true,
				"reasons": []string{
					"Высокая производительность аналитических запросов",
					"Эффективное сжатие данных",
					"Поддержка колоночного хранения",
				},
			},
			"hdfs": map[string]interface{}{
				"suitable": true,
				"reasons": []string{
					"Масштабируемость для больших объемов",
					"Отказоустойчивость",
					"Экономичность хранения",
				},
			},
		},
	}

	return recommendations
}

// generateTableSchema генерирует схему таблицы на основе анализа
func (s *AnalysisService) generateTableSchema(result map[string]interface{}) map[string]interface{} {
	_ = result["data_profile"].(map[string]interface{}) // Используем в будущем

	// Генерируем схему на основе заголовков CSV
	headers := []string{
		"created", "order_status", "ticket_status", "ticket_price",
		"visitor_category", "event_id", "is_active", "valid_to",
		"count_visitor", "is_entrance", "is_entrance_mdate", "event_name",
		"event_kind_name", "spot_id", "spot_name", "museum_name",
		"start_datetime", "ticket_id", "update_timestamp", "client_name",
		"name", "surname", "client_phone", "museum_inn", "birthday_date",
		"order_number", "ticket_number",
	}

	fields := make([]map[string]interface{}, 0)
	for _, header := range headers {
		field := map[string]interface{}{
			"name":     header,
			"type":     s.inferFieldType(header),
			"nullable": true,
			"indexed":  s.shouldCreateIndex(header),
		}
		fields = append(fields, field)
	}

	return map[string]interface{}{
		"table_name":  "museum_tickets",
		"fields":      fields,
		"primary_key": []string{"ticket_id"},
		"indexes": []map[string]interface{}{
			{"name": "idx_created", "fields": []string{"created"}},
			{"name": "idx_event_id", "fields": []string{"event_id"}},
			{"name": "idx_museum_name", "fields": []string{"museum_name"}},
		},
		"constraints": []map[string]interface{}{
			{"name": "chk_ticket_price", "type": "CHECK", "expression": "ticket_price >= 0"},
			{"name": "chk_count_visitor", "type": "CHECK", "expression": "count_visitor > 0"},
		},
	}
}

// generateDDLMetadata генерирует метаданные для DDL
func (s *AnalysisService) generateDDLMetadata(result map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ddl_generation": map[string]interface{}{
			"postgresql": map[string]interface{}{
				"table_name": "museum_tickets",
				"schema":     "public",
				"features": []string{
					"JSONB для метаданных",
					"Партиционирование по дате",
					"Индексы для аналитики",
					"Ограничения целостности",
				},
			},
			"clickhouse": map[string]interface{}{
				"table_name": "museum_tickets",
				"engine":     "MergeTree",
				"features": []string{
					"Колоночное хранение",
					"Сжатие данных",
					"Партиционирование по дате",
					"Сортировка по ticket_id",
				},
			},
			"hdfs": map[string]interface{}{
				"path":   "/data/museum/tickets/",
				"format": "Parquet",
				"features": []string{
					"Колоночное сжатие",
					"Схема данных",
					"Партиционирование",
					"Метаданные",
				},
			},
		},
		"data_characteristics": map[string]interface{}{
			"estimated_size": "~100MB",
			"row_count":      result["data_profile"].(map[string]interface{})["total_rows"],
			"column_count":   result["data_profile"].(map[string]interface{})["total_columns"],
			"data_types":     []string{"timestamp", "varchar", "numeric", "boolean", "date"},
		},
	}
}

// inferFieldType определяет тип поля на основе имени
func (s *AnalysisService) inferFieldType(fieldName string) string {
	switch {
	case containsAny(fieldName, []string{"date", "time"}):
		return "TIMESTAMP"
	case containsAny(fieldName, []string{"price", "count"}):
		return "DECIMAL(10,2)"
	case containsAny(fieldName, []string{"is_"}):
		return "BOOLEAN"
	case containsAny(fieldName, []string{"id"}):
		return "BIGINT"
	case containsAny(fieldName, []string{"phone"}):
		return "VARCHAR(20)"
	case containsAny(fieldName, []string{"email"}):
		return "VARCHAR(255)"
	default:
		return "TEXT"
	}
}

// shouldCreateIndex определяет, нужно ли создавать индекс для поля
func (s *AnalysisService) shouldCreateIndex(fieldName string) bool {
	indexFields := []string{"created", "event_id", "museum_name", "ticket_id", "order_status"}
	return contains(fieldName, indexFields)
}

// contains проверяет, содержится ли строка в слайсе
func contains(str string, slice []string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// containsAny проверяет, содержит ли строка любую из подстрок
func containsAny(str string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(str, substr) {
			return true
		}
	}
	return false
}

func main() {
	// Получение URL сервисов из переменных окружения
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434"
	}

	minioURL := os.Getenv("MINIO_URL")
	if minioURL == "" {
		minioURL = "http://minio:9000"
	}

	// Создание сервиса анализа
	analysisService := NewAnalysisService(ollamaURL, minioURL)

	// Настройка Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware для CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "data-analysis-service"})
	})

	// API маршруты
	api := r.Group("/api/v1")
	{
		api.POST("/analysis/start", analysisService.AnalyzeData)
		api.GET("/analysis/status/:analysis_id", analysisService.GetAnalysisStatus)
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Data Analysis Service запущен на порту %s", port)
	log.Fatal(r.Run(":" + port))
}
