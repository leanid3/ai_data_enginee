package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger интерфейс для логирования
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithContext(ctx context.Context) Logger
}

// zerologLogger реализация Logger с использованием zerolog
type zerologLogger struct {
	logger zerolog.Logger
}

// NewLogger создает новый логгер
func NewLogger(level, format, output string) Logger {
	// Устанавливаем уровень логирования
	zerologLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		zerologLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(zerologLevel)

	// Настраиваем формат вывода
	var logger zerolog.Logger
	if format == "json" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Pretty logging для разработки
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	// Настраиваем глобальный логгер
	log.Logger = logger

	return &zerologLogger{logger: logger}
}

// Debug логирует сообщение уровня DEBUG
func (l *zerologLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Info логирует сообщение уровня INFO
func (l *zerologLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Warn логирует сообщение уровня WARN
func (l *zerologLogger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Error логирует сообщение уровня ERROR
func (l *zerologLogger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Fatal логирует сообщение уровня FATAL и завершает программу
func (l *zerologLogger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Debugf логирует форматированное сообщение уровня DEBUG
func (l *zerologLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

// Infof логирует форматированное сообщение уровня INFO
func (l *zerologLogger) Infof(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

// Warnf логирует форматированное сообщение уровня WARN
func (l *zerologLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

// Errorf логирует форматированное сообщение уровня ERROR
func (l *zerologLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}

// Fatalf логирует форматированное сообщение уровня FATAL и завершает программу
func (l *zerologLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Msgf(format, v...)
}

// WithField добавляет поле к логгеру
func (l *zerologLogger) WithField(key string, value interface{}) Logger {
	return &zerologLogger{logger: l.logger.With().Interface(key, value).Logger()}
}

// WithFields добавляет поля к логгеру
func (l *zerologLogger) WithFields(fields map[string]interface{}) Logger {
	logger := l.logger
	for key, value := range fields {
		logger = logger.With().Interface(key, value).Logger()
	}
	return &zerologLogger{logger: logger}
}

// WithContext добавляет контекст к логгеру
func (l *zerologLogger) WithContext(ctx context.Context) Logger {
	return &zerologLogger{logger: l.logger.With().Ctx(ctx).Logger()}
}

// GetRequestID извлекает request_id из контекста
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// WithRequestID добавляет request_id к логгеру
func WithRequestID(logger Logger, requestID string) Logger {
	return logger.WithField("request_id", requestID)
}

// LogError логирует ошибку с дополнительной информацией
func LogError(logger Logger, err error, msg string, fields ...map[string]interface{}) {
	loggerWithFields := logger
	for _, fieldMap := range fields {
		loggerWithFields = loggerWithFields.WithFields(fieldMap)
	}
	loggerWithFields.WithField("error", err.Error()).Error(msg)
}

// GetLoggerFromContext извлекает логгер из контекста
func GetLoggerFromContext(ctx context.Context) Logger {
	if log, ok := ctx.Value("logger").(Logger); ok {
		return log
	}
	// Возвращаем глобальный логгер если не найден в контексте
	return NewLogger("info", "json", "stdout")
}
