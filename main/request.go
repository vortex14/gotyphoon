package main

import (
	"context"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/data/fake"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/task"
)

func init()  {
	log.InitD()
}

func main() {

	ctxGroup := task.NewTaskCtx(fake.CreateDefaultTask())

	(&forms.PipelineGroup{
		BaseLabel: interfaces.BaseLabel{
			Name:        "Http strategy",
			Required:    true,
		},
		Stages: []interfaces.BasePipelineInterface{
			&pipelines.TaskPipeline{
				BasePipeline: &forms.BasePipeline{
					Name: "prepare request",
				},
				Fn: func(context context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface) (error, context.Context) {

					newCtx := ctx.Update(context, "key", "CONTEXT DATA VALUE")


					return nil, newCtx
				},
			},
			&forms.BasePipeline{
				Name: "Request",
				Fn: func(context context.Context, logger interfaces.LoggerInterface) (error, context.Context) {

					return nil, context
				},
			},
			&forms.BasePipeline{
				Name: "SECOND STEP 2",
				Fn: func(context context.Context, logger interfaces.LoggerInterface) (error, context.Context) {


					return nil, context
				},
			},
			&pipelines.TaskPipeline{
				BasePipeline: &forms.BasePipeline{
					Name: "task-pipeline",
				},
				Fn: func(context context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface) (error, context.Context) {
					ctxData := ctx.Get(context, "key")

					logger.Info(ctxData)

					return nil, context
				},
			},
		},
		Consumers:   nil,
	}).Run(ctxGroup)
}
