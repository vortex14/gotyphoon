package interfaces

import "go.uber.org/zap"

const (
	DEBUG   = "DEBUG"
	WARNING = "WARNING"
	INFO    = "INFO"
	ERROR   = "ERROR"
)

type LoggerInterface interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
}
