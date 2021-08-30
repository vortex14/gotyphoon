package main

import (
	"context"
	"github.com/vortex14/gotyphoon/data/fake"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/logger"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)



func init() {
	(&logger.TyphoonLogger{
		Name: "App",
		Options: logger.Options{
			BaseLoggerOptions: &interfaces.BaseLoggerOptions{
				Name:          "Test-App",
				Level:         "DEBUG",
				ShowLine:      true,
				ShowFile:      true,
				ShortFileName: true,
				FullTimestamp: true,
			},
		},
	}).Init()
}


func main() {


	fakeTask, _ := fake.CreateFakeTask(interfaces.FakeTaskOptions{
		UserAgent:   false,
		Cookies:     false,
		Auth:        false,
		Proxy:       false,
		AllowedHttp: nil,
	})


	ctxGroup := context.WithValue(context.Background(), interfaces.ContextKey(interfaces.TASK), fakeTask)

	(&forms.PipelineGroup{
		BaseLabel: interfaces.BaseLabel{
			Name:        "BASE-GROUP",
			Required:    true,
		},
		Stages: []interfaces.BasePipelineInterface{
			&forms.BasePipeline{
				Name: "FIST STEP 1",
				Middlewares: []interfaces.MiddlewareInterface{
					&interfaces.Middleware{
						Required: true,
						Name:        "middleware for FIST STEP 1",
						Callback: func(context context.Context, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
							logger.Debug("run first callback for pipeline step 1")
						},
					},
					&interfaces.TaskMiddleware{
						Middleware: &interfaces.Middleware{
							Name: "task-middleware",
						},
						Callback: func(ctx context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
							//reject(Errors.ForceSkipMiddlewares)

							newCtx := interfaces.UpdateContext(ctx, "test-ctx-key", "test-ctx-data-1")

							//ctxData := interfaces.GetContextValue(newCtx, "test-ctx-key")
							//println("CTX DATA", ctxData)
							//reject(Errors.ForceSkipMiddlewares)
							next(newCtx)

						},
					},
					&interfaces.TaskMiddleware{
						Middleware: &interfaces.Middleware{
							Name: "task-middleware-2",
						},
						Callback: func(ctx context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
							ctxData := interfaces.GetContextValue(ctx, "test-ctx-key")
							//println("CTX DATA", ctxData)
							logger.Info("not skip !, new CONTEXT ::: ", ctxData)
						},
					},
				},
				LambdaHandler: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
					ctxData := interfaces.GetContextValue(ctx, "test-ctx-key")
					if ctxData  != nil{
						logger.Info("CTX DATA:", ctxData.(string))
					} else {
						logger.Error("Not found CTX DATA")
					}

					return nil, ctx
				},
			},
			&forms.BasePipeline{
				Name: "SECOND STEP 2",
				LambdaHandler: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {

					return nil, ctx
				},
			},
			&forms.TaskPipeline{
				BasePipeline: &forms.BasePipeline{
					Name: "task-pipeline",
				},
				TaskHandler: func(ctx context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface) (error, context.Context) {
					ctxData := interfaces.GetContextValue(ctx, "test-ctx-key")
					if ctxData  != nil{
						logger.Info("FOUND CTX DATA: in LambdaHandler", ctxData.(string))
					} else {
						logger.Error("Not found CTX DATA !!!")
					}

					return nil, ctx
				},
			},
		},
		Consumers:   nil,
	}).Run(ctxGroup)
}
