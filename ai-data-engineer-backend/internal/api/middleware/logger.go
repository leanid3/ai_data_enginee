package middleware

import (
	"context"
	"time"

	"ai-data-engineer-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger middleware для логирования HTTP запросов
func RequestLogger(log logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Создаем структурированный лог
		fields := map[string]interface{}{
			"timestamp":  param.TimeStamp.Format(time.RFC3339),
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"user_agent": param.Request.UserAgent(),
		}

		// Добавляем request_id если есть
		if requestID := param.Keys["request_id"]; requestID != nil {
			fields["request_id"] = requestID
		}

		// Логируем в зависимости от статуса
		if param.StatusCode >= 500 {
			log.WithFields(fields).Error("HTTP request error")
		} else if param.StatusCode >= 400 {
			log.WithFields(fields).Warn("HTTP request warning")
		} else {
			log.WithFields(fields).Info("HTTP request")
		}

		return ""
	})
}

// RequestID middleware для добавления request_id
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Добавляем request_id в контекст
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)

		// Добавляем в заголовки ответа
		c.Header("X-Request-ID", requestID)

		// Сохраняем в keys для доступа в других middleware
		c.Set("request_id", requestID)

		c.Next()
	}
}

// LoggingContext middleware для добавления логгера в контекст
func LoggingContext(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем request_id из контекста
		requestID := GetRequestID(c.Request.Context())

		// Создаем логгер с request_id
		requestLogger := logger.WithRequestID(log, requestID)

		// Добавляем логгер в контекст
		ctx := context.WithValue(c.Request.Context(), "logger", requestLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetLogger извлекает логгер из контекста
func GetLogger(c *gin.Context) logger.Logger {
	if log, ok := c.Request.Context().Value("logger").(logger.Logger); ok {
		return log
	}
	// Возвращаем глобальный логгер если не найден в контексте
	return logger.NewLogger("info", "json", "stdout")
}

// GetRequestID извлекает request_id из контекста
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

