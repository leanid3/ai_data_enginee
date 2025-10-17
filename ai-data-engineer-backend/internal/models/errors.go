package models

import (
	"fmt"
	"net/http"
)

// ErrorCode тип для кодов ошибок
type ErrorCode string

const (
	// Validation errors
	ErrorCodeValidation    ErrorCode = "validation_error"
	ErrorCodeInvalidInput  ErrorCode = "invalid_input"
	ErrorCodeMissingField  ErrorCode = "missing_field"
	ErrorCodeInvalidFormat ErrorCode = "invalid_format"

	// File errors
	ErrorCodeFileNotFound    ErrorCode = "file_not_found"
	ErrorCodeFileTooLarge    ErrorCode = "file_too_large"
	ErrorCodeUnsupportedType ErrorCode = "unsupported_file_type"
	ErrorCodeUploadFailed    ErrorCode = "upload_failed"

	// Database errors
	ErrorCodeDatabaseError     ErrorCode = "database_error"
	ErrorCodeConnectionFailed  ErrorCode = "connection_failed"
	ErrorCodeQueryFailed       ErrorCode = "query_failed"
	ErrorCodeTransactionFailed ErrorCode = "transaction_failed"

	// LLM errors
	ErrorCodeLLMError           ErrorCode = "llm_error"
	ErrorCodeLLMTimeout         ErrorCode = "llm_timeout"
	ErrorCodeLLMUnavailable     ErrorCode = "llm_unavailable"
	ErrorCodeLLMInvalidResponse ErrorCode = "llm_invalid_response"

	// Pipeline errors
	ErrorCodePipelineNotFound ErrorCode = "pipeline_not_found"
	ErrorCodePipelineFailed   ErrorCode = "pipeline_failed"
	ErrorCodeExecutionFailed  ErrorCode = "execution_failed"

	// Storage errors
	ErrorCodeStorageError       ErrorCode = "storage_error"
	ErrorCodeStorageUnavailable ErrorCode = "storage_unavailable"
	ErrorCodeBucketNotFound     ErrorCode = "bucket_not_found"

	// General errors
	ErrorCodeInternalError      ErrorCode = "internal_error"
	ErrorCodeServiceUnavailable ErrorCode = "service_unavailable"
	ErrorCodeUnauthorized       ErrorCode = "unauthorized"
	ErrorCodeForbidden          ErrorCode = "forbidden"
	ErrorCodeNotFound           ErrorCode = "not_found"
	ErrorCodeConflict           ErrorCode = "conflict"
)

// AppError представляет ошибку приложения
type AppError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	HTTPCode  int                    `json:"-"`
	Cause     error                  `json:"-"`
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap возвращает причину ошибки
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError создает новую ошибку приложения
func NewAppError(code ErrorCode, message string, httpCode int) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
	}
}

// NewAppErrorWithCause создает новую ошибку с причиной
func NewAppErrorWithCause(code ErrorCode, message string, httpCode int, cause error) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
		Cause:    cause,
	}
}

// NewValidationError создает ошибку валидации
func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:     ErrorCodeValidation,
		Message:  message,
		HTTPCode: http.StatusBadRequest,
		Details:  details,
	}
}

// NewFileNotFoundError создает ошибку "файл не найден"
func NewFileNotFoundError(fileID string) *AppError {
	return &AppError{
		Code:     ErrorCodeFileNotFound,
		Message:  "Файл не найден",
		HTTPCode: http.StatusNotFound,
		Details:  map[string]interface{}{"file_id": fileID},
	}
}

// NewDatabaseError создает ошибку базы данных
func NewDatabaseError(message string, cause error) *AppError {
	return &AppError{
		Code:     ErrorCodeDatabaseError,
		Message:  message,
		HTTPCode: http.StatusInternalServerError,
		Cause:    cause,
	}
}

// NewLLMError создает ошибку LLM
func NewLLMError(message string, cause error) *AppError {
	return &AppError{
		Code:     ErrorCodeLLMError,
		Message:  message,
		HTTPCode: http.StatusBadGateway,
		Cause:    cause,
	}
}

// NewPipelineNotFoundError создает ошибку "пайплайн не найден"
func NewPipelineNotFoundError(pipelineID string) *AppError {
	return &AppError{
		Code:     ErrorCodePipelineNotFound,
		Message:  "Пайплайн не найден",
		HTTPCode: http.StatusNotFound,
		Details:  map[string]interface{}{"pipeline_id": pipelineID},
	}
}

// NewInternalError создает внутреннюю ошибку
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Code:     ErrorCodeInternalError,
		Message:  message,
		HTTPCode: http.StatusInternalServerError,
		Cause:    cause,
	}
}

// IsAppError проверяет, является ли ошибка AppError
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

