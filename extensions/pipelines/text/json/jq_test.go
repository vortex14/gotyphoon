package json

import (
	"context"
	Errors "errors"
	"github.com/itchyny/gojq"
	"github.com/vortex14/gotyphoon/elements/forms"
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

func TestCheckJSON(t *testing.T) {

	Convey("test jq decode", t, func() {

		var cancelErr error

		var err error

		var result string

		pg := &forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{
				Name:     "Skip-Pipeline",
				Required: true,
			},
			Stages: []interfaces.BasePipelineInterface{
				&ResponseJQPipeline{
					Settings: JQSettings{Query: ".s"},
					BasePipeline: &forms.BasePipeline{
						Options: forms.GetNotRetribleOptions(),
						MetaInfo: &label.MetaInfo{
							Name:     "Skip-Pipeline",
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
						logger interfaces.LoggerInterface,
						request *http.Request,
						response *http.Response,
						data *string, jq gojq.Iter) (error, context.Context) {

						v, ok := jq.Next()

						if !ok {
							return Errors.New("not found"), context
						}

						if err, ok = v.(error); ok {
							return Errors.New("not found"), context
						}

						result = v.(string)

						return nil, context
					},
					Cn: func(err error, context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {
						logger.Error(err)
						cancelErr = err
					},
				},
			},
		}

		err = pg.Run(context.Background())

		So(err, ShouldBeNil)
		So(cancelErr, ShouldBeNil)

		So(result, ShouldEqual, "test_json")

	})
}
