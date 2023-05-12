package json

import (
	"context"
	Errors "errors"
	"fmt"
	"github.com/itchyny/gojq"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
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

func TestGenerateJQRule(t *testing.T) {
	Convey("test1", t, func() {
		topicsPrefixes := []string{
			"keepa_priority",
			"upcitemdb_priority",
			"target-com",
			"bestbuy-com-api",
		}

		var jqRule = "jq '.topics | .[] | select("

		for i, _prefix := range topicsPrefixes {
			jqRule += fmt.Sprintf(`(.topic_name | test("%s"))`, _prefix)

			if i != len(topicsPrefixes)-1 {
				jqRule += " or "
			} else {
				jqRule += ")'"
			}
		}

		println(jqRule)
	})
}

func TestNSQParseJSON(t *testing.T) {

	Convey("test", t, func() {

		var cancelErr error

		var err error

		var result string

		pg := &forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{
				Name:     "main group",
				Required: true,
			},
			Stages: []interfaces.BasePipelineInterface{
				&forms.BasePipeline{
					MetaInfo: &label.MetaInfo{Name: "prepare"},
					Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
						_task := fake.CreateDefaultTask()

						_task.SetFetcherUrl("http://typhoon-s1.ru:4151/stats?format=json")
						_task.SetFetcherTimeout(100)

						ctx = task.NewTaskCtx(_task)
						return nil, ctx

					},
				},
				net_http.CreateRequestPipeline(),
				&ResponseJQPipeline{
					Settings: JQSettings{Query: ".topics"},
					BasePipeline: &forms.BasePipeline{
						Options: forms.GetNotRetribleOptions(),
						MetaInfo: &label.MetaInfo{
							Name:     "parse JSON",
							Required: true,
						},
					},
					Fn: func(context context.Context,
						task interfaces.TaskInterface,
						logger interfaces.LoggerInterface,
						request *http.Request,
						response *http.Response,
						data *string, jq gojq.Iter) (error, context.Context) {

						//logger.Error()

						v, ok := jq.Next()
						//
						if !ok {
							return Errors.New("not found"), context
						}

						logger.Debugf("%+v", v)
						//
						//if err, ok = v.(error); ok {
						//	return Errors.New("not found"), context
						//}
						//
						//result = v.(string)

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
