package log

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/proto"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
)

// NewClientZapLogger alias for NewZapLogger with specific client ID.
func NewClientZapLogger(logLevel string, clientID string) *zap.Logger {
	return NewZapLogger(logLevel).With(zap.String("client-id", clientID))
}

var (
	// JSONPBMarshaler is the marshaller used for serializing protobuf messages.
	// If needed, this variable can be reassigned with a different marshaller with the same Marshal() signature.
	JSONPBMarshaler = &protojson.MarshalOptions{}
)

// JSONPBObjectMarshaler ...
type JSONPBObjectMarshaler struct {
	Pb proto.Message
}

// NewZapLogger returns new Zap logger.
func NewZapLogger(logLevel string) *zap.Logger {
	var level zap.AtomicLevel

	switch logLevel {
	case DebugLevel:
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case InfoLevel:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WarnLevel:
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case ErrorLevel:
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case FatalLevel:
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case PanicLevel:
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            level,
		Development:      false,
		Sampling:         nil,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	logger, err := config.Build(zap.AddCallerSkip(0))
	if err != nil {
		log.Printf("failed build zap log: %v", err)
		return zap.NewNop()
	}

	zap.RedirectStdLog(logger)

	return logger
}

// MarshalLogObject ...
func (j *JSONPBObjectMarshaler) MarshalLogObject(e zapcore.ObjectEncoder) error {
	// ZAP jsonEncoder deals with AddReflect by using json.MarshalObject. The same thing applies for consoleEncoder.
	return e.AddReflected("msg", j)
}

// MarshalJSON ...
func (j *JSONPBObjectMarshaler) MarshalJSON() ([]byte, error) {
	b, err := JSONPBMarshaler.Marshal(j.Pb)
	if err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %w", err)
	}

	return b, nil
}

// NewJSONPBObjectMarshaller init new jsonpb object marshaller.
func NewJSONPBObjectMarshaller(msg proto.Message) *JSONPBObjectMarshaler {
	return &JSONPBObjectMarshaler{Pb: msg}
}

// NewObservedZapLogger creates logger that buffers logs in memory (without any encoding).
// It's particularly useful in tests.
func NewObservedZapLogger() *zap.Logger {
	observedZapCore, _ := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	return observedLogger
}

// LogProtoMessageAsJSON  logs proto message as json.
func LogProtoMessageAsJSON(logger *zap.Logger, level zapcore.Level, pbMsg interface{}, msg string, key string) {
	if p, ok := pbMsg.(proto.Message); ok {
		logger.Check(level, msg).Write(zap.Object(key, &JSONPBObjectMarshaler{Pb: p}))
	}
}
