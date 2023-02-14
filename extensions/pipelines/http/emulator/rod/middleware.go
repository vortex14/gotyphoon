package rod

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/middlewares"
	"github.com/vortex14/gotyphoon/interfaces"
)

func ConstructorRodProxyRequestMiddleware(required bool) interfaces.MiddlewareInterface {
	return &middlewares.TaskMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required: required,
				Name:     "prepare for rod request",
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

			url := launcher.New().Proxy(task.GetProxyAddress()).Delete("use-mock-keychain").MustLaunch()
			browser := rod.New().ControlURL(url)

			browser = browser.WithPanic(func(i interface{}) {
				logger.Error(fmt.Sprintf(">>>>> %v", i))
				reject(fmt.Errorf("undefined behavior: %v", i))

			})

			err := browser.Connect()

			if err != nil {
				logger.Error("Need to retry request after later, browser is very busy or doesn't accessible ")
				reject(err)
				return
			}

			context = NewBrowserCtx(context, browser)
			next(context)

		},
	}
}

func ConstructorRodBasicRequestMiddleware(required bool) interfaces.MiddlewareInterface {
	return &middlewares.TaskMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required: required,
				Name:     "prepare for rod request",
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

			url := launcher.New().Delete("use-mock-keychain").MustLaunch()
			browser := rod.New().ControlURL(url)

			browser = browser.WithPanic(func(i interface{}) {
				logger.Error(fmt.Errorf("undefined behavior: %v", i))
			})

			err := browser.Connect()

			if err != nil {
				logger.Error("Need to retry request after later, browser is very busy or doesn't accessible ")
				reject(err)
				return
			}

			context = NewBrowserCtx(context, browser)
			next(context)

		},
	}
}
