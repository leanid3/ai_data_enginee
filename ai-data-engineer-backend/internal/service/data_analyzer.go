package service

import (
	"context"
	"fmt"
	"strings"

	"ai-data-engineer-backend/internal/models"
	"ai-data-engineer-backend/pkg/logger"
)

// dataAnalyzer реализация DataAnalyzer
type dataAnalyzer struct {
	logger logger.Logger
}

// NewDataAnalyzer создает новый анализатор данных
func NewDataAnalyzer(logger logger.Logger) DataAnalyzer {
	return &dataAnalyzer{
		logger: logger,
	}
}

// AnalyzeDataStructure анализирует структуру данных
func (a *dataAnalyzer) AnalyzeDataStructure(ctx context.Context, profile *models.DataProfile) (*models.AnalysisResult, error) {
	a.logger.Info("Analyzing data structure")

	// Анализируем качество данных
	qualityScore := a.CalculateDataQuality(ctx, profile)
	_ = qualityScore // Используем переменную

	// Генерируем рекомендации
	recommendations := a.GenerateRecommendations(ctx, profile)

	// Определяем рекомендуемое хранилище
	storageRecommendation := a.determineStorageRecommendation(profile)

	// Создаем схему таблицы
	tableSchema := a.generateTableSchema(profile)

	// Создаем DDL метаданные
	ddlMetadata := a.generateDDLMetadata(profile)

	return &models.AnalysisResult{
		AnalysisID:            fmt.Sprintf("analysis_%d", profile.CreatedAt.Unix()),
		DataProfile:           *profile,
		Recommendations:       recommendations,
		StorageRecommendation: storageRecommendation,
		TableSchema:           tableSchema,
		DDLMetadata:           ddlMetadata,
		Status:                models.AnalysisStatusCompleted,
	}, nil
}

// CalculateDataQuality рассчитывает качество данных
func (a *dataAnalyzer) CalculateDataQuality(ctx context.Context, profile *models.DataProfile) float64 {
	a.logger.Info("Calculating data quality")

	// Базовое качество на основе профиля
	quality := profile.DataQualityScore

	// Дополнительные факторы
	if profile.TotalRows > 1000 {
		quality += 0.1 // Больше данных = лучше
	}

	if len(profile.Fields) > 5 {
		quality += 0.05 // Больше полей = сложнее, но интереснее
	}

	// Проверяем наличие обязательных полей
	hasID := false
	hasTimestamp := false
	for _, field := range profile.Fields {
		if strings.Contains(strings.ToLower(field.Name), "id") {
			hasID = true
		}
		if strings.Contains(strings.ToLower(field.Name), "time") || strings.Contains(strings.ToLower(field.Name), "date") {
			hasTimestamp = true
		}
	}

	if hasID {
		quality += 0.1
	}
	if hasTimestamp {
		quality += 0.05
	}

	// Ограничиваем качество от 0 до 1
	if quality > 1.0 {
		quality = 1.0
	}
	if quality < 0.0 {
		quality = 0.0
	}

	return quality
}

// GenerateRecommendations генерирует рекомендации
func (a *dataAnalyzer) GenerateRecommendations(ctx context.Context, profile *models.DataProfile) []string {
	a.logger.Info("Generating recommendations")

	recommendations := []string{}

	// Рекомендации по качеству данных
	if profile.DataQualityScore < 0.7 {
		recommendations = append(recommendations, "Низкое качество данных - рекомендуется очистка и валидация")
	} else if profile.DataQualityScore < 0.9 {
		recommendations = append(recommendations, "Среднее качество данных - рекомендуется проверка на аномалии")
	} else {
		recommendations = append(recommendations, "Высокое качество данных - готовы к использованию")
	}

	// Рекомендации по объему данных
	if profile.TotalRows > 1000000 {
		recommendations = append(recommendations, "Большой объем данных - рекомендуется партиционирование")
	} else if profile.TotalRows > 100000 {
		recommendations = append(recommendations, "Средний объем данных - рекомендуется индексирование")
	}

	// Рекомендации по типу файла
	switch profile.DataType {
	case "csv":
		recommendations = append(recommendations, "CSV формат - рекомендуется валидация разделителей")
	case "json":
		recommendations = append(recommendations, "JSON формат - рекомендуется проверка структуры")
	case "xml":
		recommendations = append(recommendations, "XML формат - рекомендуется валидация схемы")
	}

	// Рекомендации по полям
	hasNumericFields := false
	hasTextFields := false
	for _, field := range profile.Fields {
		if field.Type == "numeric" {
			hasNumericFields = true
		}
		if field.Type == "string" {
			hasTextFields = true
		}
	}

	if hasNumericFields {
		recommendations = append(recommendations, "Найдены числовые поля - рекомендуется статистический анализ")
	}
	if hasTextFields {
		recommendations = append(recommendations, "Найдены текстовые поля - рекомендуется анализ текста")
	}

	return recommendations
}

// determineStorageRecommendation определяет рекомендуемое хранилище
func (a *dataAnalyzer) determineStorageRecommendation(profile *models.DataProfile) models.StorageRecommendation {
	// Анализируем характеристики данных
	fileType := profile.DataType
	totalRows := profile.TotalRows
	qualityScore := profile.DataQualityScore

	// Определяем основное хранилище
	var primaryStorage string
	var secondaryStorage []string

	if totalRows > 1000000 {
		// Большие объемы - ClickHouse или HDFS
		primaryStorage = "ClickHouse"
		secondaryStorage = []string{"HDFS", "PostgreSQL"}
	} else if qualityScore > 0.8 && fileType == "csv" {
		// Структурированные данные с хорошим качеством - PostgreSQL
		primaryStorage = "PostgreSQL"
		secondaryStorage = []string{"ClickHouse", "HDFS"}
	} else {
		// По умолчанию - PostgreSQL
		primaryStorage = "PostgreSQL"
		secondaryStorage = []string{"ClickHouse", "HDFS"}
	}

	reasoning := map[string]interface{}{
		"file_type":     fileType,
		"data_volume":   totalRows,
		"quality_score": qualityScore,
		"recommendation": fmt.Sprintf("%s подходит для данных с объемом %d строк и качеством %.2f",
			primaryStorage, totalRows, qualityScore),
	}

	storageOptions := map[string]interface{}{
		"postgresql": map[string]interface{}{
			"suitable": primaryStorage == "PostgreSQL",
			"reasons": []string{
				"Структурированные данные",
				"Хорошее качество данных",
				"Поддержка ACID транзакций",
			},
		},
		"clickhouse": map[string]interface{}{
			"suitable": primaryStorage == "ClickHouse",
			"reasons": []string{
				"Высокая производительность аналитических запросов",
				"Эффективное сжатие данных",
			},
		},
		"hdfs": map[string]interface{}{
			"suitable": primaryStorage == "HDFS",
			"reasons": []string{
				"Масштабируемость для больших объемов",
				"Отказоустойчивость",
			},
		},
	}

	return models.StorageRecommendation{
		PrimaryStorage:   primaryStorage,
		SecondaryStorage: secondaryStorage,
		Reasoning:        reasoning,
		StorageOptions:   storageOptions,
	}
}

// generateTableSchema генерирует схему таблицы
func (a *dataAnalyzer) generateTableSchema(profile *models.DataProfile) models.TableSchema {
	fields := make([]models.TableField, len(profile.Fields))
	indexes := []models.TableIndex{}
	constraints := []models.TableConstraint{}

	for i, field := range profile.Fields {
		tableField := models.TableField{
			Name:        field.Name,
			Type:        a.mapToSQLType(field.Type),
			Nullable:    field.Nullable,
			Indexed:     a.shouldCreateIndex(field.Name),
			Description: field.Description,
		}
		fields[i] = tableField

		// Создаем индекс для важных полей
		if tableField.Indexed {
			indexes = append(indexes, models.TableIndex{
				Name:   fmt.Sprintf("idx_%s", field.Name),
				Fields: []string{field.Name},
				Type:   "btree",
			})
		}
	}

	// Добавляем ограничения
	if len(fields) > 0 {
		constraints = append(constraints, models.TableConstraint{
			Name:       "chk_data_quality",
			Type:       "CHECK",
			Expression: "data_quality_score > 0",
		})
	}

	return models.TableSchema{
		TableName:   "analyzed_data",
		Fields:      fields,
		PrimaryKey:  []string{"id"},
		Indexes:     indexes,
		Constraints: constraints,
	}
}

// generateDDLMetadata генерирует метаданные DDL
func (a *dataAnalyzer) generateDDLMetadata(profile *models.DataProfile) models.DDLMetadata {
	return models.DDLMetadata{
		DDLGeneration: map[string]interface{}{
			"postgresql": map[string]interface{}{
				"table_name": "analyzed_data",
				"schema":     "public",
				"features": []string{
					"JSONB для метаданных",
					"Партиционирование по дате",
					"Индексы для аналитики",
				},
			},
			"clickhouse": map[string]interface{}{
				"table_name": "analyzed_data",
				"engine":     "MergeTree",
				"features": []string{
					"Колоночное хранение",
					"Сжатие данных",
					"Партиционирование",
				},
			},
		},
		DataCharacteristics: map[string]interface{}{
			"estimated_size": fmt.Sprintf("~%dMB", profile.FileSize/1024/1024),
			"row_count":      profile.TotalRows,
			"column_count":   len(profile.Fields),
			"data_types":     a.getDataTypes(profile.Fields),
		},
	}
}

// mapToSQLType маппит типы данных в SQL
func (a *dataAnalyzer) mapToSQLType(dataType string) string {
	switch dataType {
	case "numeric":
		return "DECIMAL(10,2)"
	case "string":
		return "VARCHAR(255)"
	case "timestamp":
		return "TIMESTAMP"
	case "boolean":
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}

// shouldCreateIndex определяет, нужно ли создавать индекс
func (a *dataAnalyzer) shouldCreateIndex(fieldName string) bool {
	indexFields := []string{"id", "created", "updated", "timestamp", "date"}
	fieldLower := strings.ToLower(fieldName)

	for _, indexField := range indexFields {
		if strings.Contains(fieldLower, indexField) {
			return true
		}
	}
	return false
}

// getDataTypes возвращает список типов данных
func (a *dataAnalyzer) getDataTypes(fields []models.DataField) []string {
	typeSet := make(map[string]bool)
	for _, field := range fields {
		typeSet[field.Type] = true
	}

	types := make([]string, 0, len(typeSet))
	for dataType := range typeSet {
		types = append(types, dataType)
	}
	return types
}
