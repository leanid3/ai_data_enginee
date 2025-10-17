package middleware

import (
	"net/http"
	"time"

	"ai-data-engineer-backend/internal/models"
	"ai-data-engineer-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware для обработки ошибок
func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Проверяем, есть ли ошибки
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Получаем request_id
			requestID := GetRequestID(c.Request.Context())

			// Логируем ошибку
			log.WithField("request_id", requestID).
				WithField("error", err.Error()).
				Error("Request error")

			// Проверяем тип ошибки
			if appErr, ok := models.IsAppError(err.Err); ok {
				// Обрабатываем AppError
				response := models.ErrorResponse{
					Error:     string(appErr.Code),
					Message:   appErr.Message,
					Details:   appErr.Details,
					RequestID: requestID,
					Timestamp: time.Now(),
				}
				c.JSON(appErr.HTTPCode, response)
			} else {
				// Обрабатываем общие ошибки
				response := models.ErrorResponse{
					Error:     "internal_error",
					Message:   "Внутренняя ошибка сервера",
					RequestID: requestID,
					Timestamp: time.Now(),
				}
				c.JSON(http.StatusInternalServerError, response)
			}
		}
	}
}

// Recovery middleware для восстановления после panic
func Recovery(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Получаем request_id
		requestID := GetRequestID(c.Request.Context())

		// Логируем panic
		log.WithField("request_id", requestID).
			WithField("panic", recovered).
			Error("Panic recovered")

		// Отправляем ответ об ошибке
		response := models.ErrorResponse{
			Error:     "internal_error",
			Message:   "Внутренняя ошибка сервера",
			RequestID: requestID,
			Timestamp: time.Now(),
		}
		c.JSON(http.StatusInternalServerError, response)
	})
}

// NotFoundHandler обработчик для 404 ошибок
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c.Request.Context())

		response := models.ErrorResponse{
			Error:     "not_found",
			Message:   "Ресурс не найден",
			RequestID: requestID,
			Timestamp: time.Now(),
		}
		c.JSON(http.StatusNotFound, response)
	}
}

// MethodNotAllowedHandler обработчик для 405 ошибок
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c.Request.Context())

		response := models.ErrorResponse{
			Error:     "method_not_allowed",
			Message:   "Метод не разрешен",
			RequestID: requestID,
			Timestamp: time.Now(),
		}
		c.JSON(http.StatusMethodNotAllowed, response)
	}
}

