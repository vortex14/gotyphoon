package html

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/errors"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"net/http"
	"testing"

	"github.com/vortex14/gotyphoon/elements/models/label"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func TestCancelCallback(t *testing.T) {

	Convey("skip stages", t, func() {

		var cancelErr error

		var err error

		pg := &forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{
				Name:     "first group",
				Required: true,
			},
			Stages: []interfaces.BasePipelineInterface{
				&ResponseHtmlPipeline{
					BasePipeline: &forms.BasePipeline{
						Options: forms.GetNotRetribleOptions(),
						MetaInfo: &label.MetaInfo{
							Name:     "html goQuery",
							Required: true,
						},
						Middlewares: []interfaces.MiddlewareInterface{
							net_http.ConstructorMockTaskMiddleware(true),
							net_http.ConstructorPrepareRequestMiddleware(true),
							net_http.ConstructorMockResponseMiddleware(true),
						},
					},
					Fn: func(context context.Context,
						task interfaces.TaskInterface,
						logger interfaces.LoggerInterface, request *http.Request, response *http.Response,
						data *string, doc *goquery.Document) (error, context.Context) {

						return errors.PipelineFailed, context
					},
					Cn: func(err error, context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {
						cancelErr = err
					},
				},
			},
		}

		err = pg.Run(context.Background())

		So(err, ShouldBeError)
		So(cancelErr, ShouldBeError)
	})
}
