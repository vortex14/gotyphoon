package rod

import (
	"context"
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
			err := browser.Connect()

			if err != nil {
				reject(err)
				return
			}

			context = NewBrowserCtx(context, browser)

			next(context)

		},
	}
}

//func ConstructorRodProxyRequestMiddleware(required bool) interfaces.MiddlewareInterface {
//	return &middlewares.TaskMiddleware{
//		Middleware: &forms.Middleware{
//			MetaInfo: &label.MetaInfo{
//				Required: required,
//				Name:     "set launcher for rod request",
//			},
//		},
//		Fn: func(context context.Context, task *task.TyphoonTask,
//			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
//
//			sB, browser := GetBrowserCtx(context)
//
//			if !sB {
//				reject(Errors.New("rod browser not found"))
//				return
//			}
//
//			url := launcher.New().Proxy(task.GetProxyAddress()).Delete("use-mock-keychain").MustLaunch()
//			browser.ControlURL(url)
//
//			next(context)
//
//		},
//	}
//}
