package main

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	httpMiddlewares "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"net/http"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/data/fake"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/middlewares"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func main20() {


	fakeTask, _ := fake.CreateFakeTask(interfaces.FakeTaskOptions{
		UserAgent:   false,
		Cookies:     false,
		Auth:        false,
		Proxy:       false,
		AllowedHttp: nil,
	})


	ctxGroup := context.WithValue(context.Background(), ctx.ContextKey(interfaces.TASK), fakeTask)

	(&forms.PipelineGroup{
		BaseLabel: interfaces.BaseLabel{
			Name:        "BASE-GROUP",
			Required:    true,
		},
		Stages: []interfaces.BasePipelineInterface{
			&pipelines.TaskPipeline{
				BasePipeline: &forms.BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name: "task-pipeline",
					},
				},
				Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) (error, context.Context) {
					return nil, nil
				},
			},
			&forms.BasePipeline{
				MetaInfo: &label.MetaInfo{
					Name: "FIST STEP 1",
				},
				Middlewares: []interfaces.MiddlewareInterface{
					&middlewares.TaskMiddleware{
						Middleware: &forms.Middleware{
							MetaInfo: &label.MetaInfo{
								Name: "task-middleware",
							},
						},
						Fn: func(context context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
							//reject(Errors.ForceSkipMiddlewares)

							newCtx := ctx.Update(context, "test-ctx-key", "test-ctx-data-1")

							//ctxData := interfaces.GetContextValue(newCtx, "test-ctx-key")
							//println("CTX DATA", ctxData)
							//reject(Errors.ForceSkipMiddlewares)
							next(newCtx)

						},
					},
					&middlewares.TaskMiddleware{
						Middleware: &forms.Middleware{
							MetaInfo: &label.MetaInfo{
								Name: "task-middleware-2",
							},
						},
						Fn: func(context context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
							ctxData := ctx.Get(context, "test-ctx-key")
							logger.Info("not skip !, new CONTEXT ::: ", ctxData)
						},
					},
					&httpMiddlewares.HttpMiddleware{
						Middleware: &forms.Middleware{
							MetaInfo: &label.MetaInfo{
								Required: true,
								Name: "task-middleware-2",
							},
						},
						Fn: func(context context.Context, task *task.TyphoonTask, request *http.Request, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
							logger.Error("!!!!!!!!!!")
						},
					},
				},
				Fn: func(context context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
					ctxData := ctx.Get(context, "test-ctx-key")
					if ctxData  != nil{
						logger.Info("CTX DATA:", ctxData.(string))
					} else {
						logger.Error("Not found CTX DATA")
					}

					return nil, context
				},
			},
			&forms.BasePipeline{
				MetaInfo: &label.MetaInfo{
					Name: "SECOND STEP 2",
				},
				Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {

					return nil, ctx
				},
			},
			&pipelines.TaskPipeline{
				BasePipeline: &forms.BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name: "task-pipeline",
					},
				},
				Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) (error, context.Context) {
					return nil, nil
				},
			},
		},
		Consumers:   nil,
	}).Run(ctxGroup)
}
