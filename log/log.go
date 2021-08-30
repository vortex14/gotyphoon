package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/extensions/logger"
	"github.com/vortex14/gotyphoon/interfaces"
)

// D is data custom log data for logger. Logrus.Field{} is equivalent D, but shorty for import
type D map[string]interface{}


func New(values map[string]interface{}) *logrus.Entry {
	return logrus.WithFields(values)
}

func NewCtx(context context.Context, logger *logrus.Entry) context.Context {
	return ctx.Update(context, interfaces.LOGGER, logger)
}

func NewCtxValues(context context.Context, values D) context.Context {
	return ctx.Update(context, interfaces.LOGGER, New(values))
}

func Get(context context.Context) (bool, interfaces.LoggerInterface) {
	log, ok := context.Value(ctx.ContextKey(interfaces.LOGGER)).(*logrus.Entry)
	return ok, log
}

// InitD is debug logger configuration
func InitD()  {
	(&logger.TyphoonLogger{
		Name: "App",
		Options: logger.Options{
			BaseLoggerOptions: &interfaces.BaseLoggerOptions{
				Name:          "App-Debug-Logger",
				Level:         "DEBUG",
				ShowLine:      true,
				ShowFile:      true,
				ShortFileName: true,
				FullTimestamp: true,
			},
		},
	}).Init()
}