package json

import (
	"context"
	"github.com/itchyny/gojq"
	"github.com/vortex14/gotyphoon/utils"
	"net/http"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"

	Errors "github.com/vortex14/gotyphoon/errors"
	netHttp "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/interfaces"
)

type JQSettings struct {
	Query string
}

type ResponseJQPipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline

	netHttp.HttpResponsePipeline

	Settings JQSettings

	Fn func(
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		request *http.Request,
		response *http.Response,
		data *string,
		jq gojq.Iter,

	) (error, context.Context)

	Cn func(
		err error,
		context context.Context,

		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *ResponseJQPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil {
		reject(t, Errors.TaskPipelineRequiredHandler)
		return
	}

	ok, taskInstance, logger, _, request, _, response, data := t.UnpackResponse(context)

	if !ok {
		reject(t, Errors.PipelineContexFailed)
		return
	}

	t.SafeRun(func() error {

		query, err := gojq.Parse(t.Settings.Query)
		if err != nil {
			return err
		}
		var model interface{}

		err = utils.JsonLoad(&model, *data)
		if err != nil {
			return err
		}

		iter := query.Run(model)
		context = NewJQCtx(context, iter)

		err, newContext := t.Fn(context, taskInstance, logger, request, response, data, iter)
		if err != nil {
			return err
		}
		next(newContext)
		return nil

	}, func(err error) {
		reject(t, err)
		t.Cancel(context, logger, err)
	})

}

func (t *ResponseJQPipeline) Cancel(
	context context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil {
		return
	}
	ok, taskInstance, logger := t.UnpackCtx(context)
	if !ok {
		return
	}
	t.Cn(err, context, taskInstance, logger)
}
