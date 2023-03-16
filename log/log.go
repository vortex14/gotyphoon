package log

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/interfaces"
)

var logOnce sync.Once

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

func Patch(logger *logrus.Entry, values map[string]interface{}) *logrus.Entry {
	return logger.WithFields(values)
}

func PatchCtx(ctx context.Context, values map[string]interface{}) context.Context {
	s, logger := Get(ctx)
	if !s {
		logger = New(map[string]interface{}{})
	}

	return NewCtx(ctx, PatchLogI(logger, values))

}

func PatchLogI(logger interfaces.LoggerInterface, values map[string]interface{}) *logrus.Entry {
	return Patch(logger.(*logrus.Entry), values)
}

func Get(context context.Context) (bool, interfaces.LoggerInterface) {
	log, ok := ctx.Get(context, interfaces.LOGGER).(*logrus.Entry)
	if !ok {
		log = New(D{"unstable context": true})
	}
	return ok, log
}

// InitD is debug logger configuration
func InitD() {
	logOnce.Do(func() {
		(&TyphoonLogger{
			Name: "App",
			Options: Options{
				BaseOptions: &BaseOptions{
					Name:          "App-Debug-Logger",
					Level:         "DEBUG",
					ShowLine:      true,
					ShowFile:      true,
					ShortFileName: false,
					FullTimestamp: true,
				},
			},
		}).Init()
	})
}

func Init(opts *BaseOptions) {
	logOnce.Do(func() {
		(&TyphoonLogger{
			Name: "App",
			Options: Options{
				BaseOptions: opts,
			},
		}).Init()
	})
}
