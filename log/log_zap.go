package log

import (
	"context"
	"go.uber.org/zap"
	"sync"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/interfaces"
)

var logOnce sync.Once

// D is data custom log data for logger. Logrus.Field{} is equivalent D, but shorty for import
type D map[string]interface{}

func New(level string, values map[string]interface{}) *zap.Logger {
	var fields []zap.Field
	for key, value := range values {
		fields = append(fields, zap.Any(key, value))
	}
	return NewZapLogger(level).With(fields...)
}

func NewCtx(context context.Context, logger *zap.Logger) context.Context {
	return ctx.Update(context, interfaces.LOGGER, logger)
}

func NewCtxValues(context context.Context, level string, values D) context.Context {
	return ctx.Update(context, interfaces.LOGGER, New(level, values))
}

func Patch(logger *zap.Logger, values map[string]interface{}) *zap.Logger {
	var fields []zap.Field
	for key, value := range values {
		fields = append(fields, zap.Any(key, value))
	}
	return logger.With(fields...)
}

func PatchCtx(ctx context.Context, values map[string]interface{}) context.Context {
	logger := New(DebugLevel, map[string]interface{}{})
	return NewCtx(ctx, PatchLogI(logger, values))

}

func PatchLogI(logger interfaces.LoggerInterface, values map[string]interface{}) *zap.Logger {
	return Patch(logger.(*zap.Logger), values)
}

func Get(context context.Context) (bool, interfaces.LoggerInterface) {
	log, ok := ctx.Get(context, interfaces.LOGGER).(*zap.Logger)
	if !ok {
		log = New(DebugLevel, D{"unstable context": true})
	}
	return ok, log
}
