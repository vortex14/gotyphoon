package html

import (
	Context "context"
	Errors "errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Task "github.com/vortex14/gotyphoon/elements/models/task"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/interfaces"
	"net/http"
)

func MakeRequestThroughProxy(task *Task.TyphoonTask, callback net_http.ValidationCallback) error {
	ctxGroup := Task.NewTaskCtx(task)

	return (&forms.PipelineGroup{
		MetaInfo: &label.MetaInfo{
			Name:     "Http strategy",
			Required: true,
		},
		Stages: []interfaces.BasePipelineInterface{
			net_http.CreateProxyRequestPipeline(&forms.Options{Retry: forms.RetryOptions{MaxCount: 2}}),
			&ResponseHtmlPipeline{
				BasePipeline: &forms.BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name: "Response pipeline",
					},
				},
				Fn: func(context Context.Context,
					task interfaces.TaskInterface, logger interfaces.LoggerInterface,
					request *http.Request, response *http.Response,
					data *string, doc *goquery.Document) (error, Context.Context) {

					status := callback(logger, response, doc)
					if !status {
						return Errors.New("not ready "), context
					}
					return nil, context
				},
			},
		},
	}).Run(ctxGroup)
}
