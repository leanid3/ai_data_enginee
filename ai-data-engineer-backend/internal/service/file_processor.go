package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ai-data-engineer-backend/internal/models"
	"ai-data-engineer-backend/pkg/logger"
)

// fileProcessor реализация FileProcessor
type fileProcessor struct {
	logger logger.Logger
}

// NewFileProcessor создает новый файловый процессор
func NewFileProcessor(logger logger.Logger) FileProcessor {
	return &fileProcessor{
		logger: logger,
	}
}

// ParseCSV парсит CSV файл
func (p *fileProcessor) ParseCSV(ctx context.Context, content []byte) (*models.DataProfile, error) {
	p.logger.Info("Parsing CSV file")

	reader := csv.NewReader(strings.NewReader(string(content)))
	reader.Comma = ';' // Используем точку с запятой как разделитель

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("empty CSV file")
	}

	// Анализируем заголовки
	headers := records[0]
	fields := make([]models.DataField, len(headers))

	// Анализируем типы данных в каждой колонке
	for i, header := range headers {
		field := models.DataField{
			Name:     strings.TrimSpace(header),
			Type:     p.inferColumnType(records[1:], i),
			Nullable: true,
		}

		// Анализируем sample данные
		if len(records) > 1 {
			field.SampleValue = records[1][i]
		}

		fields[i] = field
	}

	// Вычисляем качество данных
	qualityScore := p.calculateDataQuality(records)

	// Получаем sample данные (первые 5 строк)
	sampleData := p.getSampleData(records, 5)

	return &models.DataProfile{
		DataType:         "csv",
		TotalRows:        len(records) - 1, // Исключаем заголовок
		SampledRows:      len(records) - 1,
		Fields:           fields,
		SampleData:       strings.Join(sampleData, "\n"),
		DataQualityScore: qualityScore,
		FileSize:         int64(len(content)),
		Encoding:         "UTF-8",
		Delimiter:        ";",
		HasHeaders:       true,
		CreatedAt:        time.Now(),
	}, nil
}

// ParseJSON парсит JSON файл
func (p *fileProcessor) ParseJSON(ctx context.Context, content []byte) (*models.DataProfile, error) {
	p.logger.Info("Parsing JSON file")

	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Анализируем структуру JSON
	fields := p.analyzeJSONStructure(data)

	return &models.DataProfile{
		DataType:         "json",
		TotalRows:        1,
		SampledRows:      1,
		Fields:           fields,
		SampleData:       string(content),
		DataQualityScore: 0.9, // JSON обычно имеет хорошее качество
		FileSize:         int64(len(content)),
		Encoding:         "UTF-8",
		HasHeaders:       false,
		CreatedAt:        time.Now(),
	}, nil
}

// ParseXML парсит XML файл
func (p *fileProcessor) ParseXML(ctx context.Context, content []byte) (*models.DataProfile, error) {
	p.logger.Info("Parsing XML file")

	var data interface{}
	if err := xml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Анализируем структуру XML
	fields := p.analyzeXMLStructure(data)

	return &models.DataProfile{
		DataType:         "xml",
		TotalRows:        1,
		SampledRows:      1,
		Fields:           fields,
		SampleData:       string(content),
		DataQualityScore: 0.8,
		FileSize:         int64(len(content)),
		Encoding:         "UTF-8",
		HasHeaders:       false,
		CreatedAt:        time.Now(),
	}, nil
}

// DetectFileType определяет тип файла по расширению
func (p *fileProcessor) DetectFileType(filename string) string {
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	switch ext {
	case "csv":
		return "csv"
	case "json":
		return "json"
	case "xml":
		return "xml"
	default:
		return "unknown"
	}
}

// inferColumnType определяет тип колонки на основе данных
func (p *fileProcessor) inferColumnType(records [][]string, columnIndex int) string {
	if len(records) == 0 {
		return "string"
	}

	// Анализируем первые несколько строк
	sampleSize := 5
	if len(records) < sampleSize {
		sampleSize = len(records)
	}

	for i := 0; i < sampleSize; i++ {
		if columnIndex >= len(records[i]) {
			continue
		}

		value := strings.TrimSpace(records[i][columnIndex])
		if value == "" {
			continue
		}

		// Проверяем, является ли значение числом
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			return "numeric"
		}

		// Проверяем, является ли значение датой
		if p.isDate(value) {
			return "timestamp"
		}

		// Проверяем, является ли значение булевым
		if p.isBoolean(value) {
			return "boolean"
		}
	}

	return "string"
}

// isDate проверяет, является ли значение датой
func (p *fileProcessor) isDate(value string) bool {
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
func (p *fileProcessor) isBoolean(value string) bool {
	valueLower := strings.ToLower(value)
	return valueLower == "true" || valueLower == "false" || valueLower == "1" || valueLower == "0"
}

// calculateDataQuality рассчитывает качество данных
func (p *fileProcessor) calculateDataQuality(records [][]string) float64 {
	if len(records) <= 1 {
		return 0.0
	}

	totalCells := 0
	emptyCells := 0

	for i := 1; i < len(records); i++ { // Пропускаем заголовок
		for j := 0; j < len(records[i]); j++ {
			totalCells++
			if strings.TrimSpace(records[i][j]) == "" {
				emptyCells++
			}
		}
	}

	if totalCells == 0 {
		return 0.0
	}

	return 1.0 - (float64(emptyCells) / float64(totalCells))
}

// getSampleData возвращает образец данных
func (p *fileProcessor) getSampleData(records [][]string, maxRows int) []string {
	if len(records) <= maxRows {
		result := make([]string, len(records))
		for i, record := range records {
			result[i] = strings.Join(record, ";")
		}
		return result
	}

	result := make([]string, maxRows)
	for i := 0; i < maxRows; i++ {
		result[i] = strings.Join(records[i], ";")
	}
	return result
}

// analyzeJSONStructure анализирует структуру JSON
func (p *fileProcessor) analyzeJSONStructure(data interface{}) []models.DataField {
	fields := []models.DataField{}

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			field := models.DataField{
				Name:     key,
				Type:     p.getJSONType(value),
				Nullable: value == nil,
			}
			fields = append(fields, field)
		}
	case []interface{}:
		if len(v) > 0 {
			// Анализируем первый элемент массива
			if obj, ok := v[0].(map[string]interface{}); ok {
				for key, value := range obj {
					field := models.DataField{
						Name:     key,
						Type:     p.getJSONType(value),
						Nullable: value == nil,
					}
					fields = append(fields, field)
				}
			}
		}
	}

	return fields
}

// analyzeXMLStructure анализирует структуру XML
func (p *fileProcessor) analyzeXMLStructure(data interface{}) []models.DataField {
	// Упрощенный анализ XML структуры
	fields := []models.DataField{
		{
			Name:     "xml_content",
			Type:     "string",
			Nullable: false,
		},
	}
	return fields
}

// getJSONType определяет тип JSON значения
func (p *fileProcessor) getJSONType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case float64:
		return "numeric"
	case bool:
		return "boolean"
	case nil:
		return "null"
	default:
		return "object"
	}
}
