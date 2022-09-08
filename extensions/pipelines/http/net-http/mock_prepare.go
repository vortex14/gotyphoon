package net_http

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Task "github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/extensions/middlewares"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"

	"context"
)

func ConstructorMockTaskMiddleware(required bool) interfaces.MiddlewareInterface {
	return &middlewares.TaskMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required: required,
				Name:     "mock task response",
			},
		},
		Fn: func(context context.Context, task *Task.TyphoonTask,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

			t := fake.CreateDefaultTask()

			l := log.New(map[string]interface{}{"test": "test"})
			ctx := log.NewCtx(context, l)

			ctx = Task.PatchCtx(ctx, t)

			next(ctx)

		},
	}
}
