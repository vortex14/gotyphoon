package rod

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Task "github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"testing"
)

func init() {
	log.InitD()
}

func TestHttpRodRequestPipeline_Run(t *testing.T) {

	Convey("Create a rod pipeline", t, func() {

		g1 := forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{
				Name: "Rod group",
			},
			Stages: []interfaces.BasePipelineInterface{
				CreateProxyRodRequestPipeline(forms.GetCustomRetryOptions(1), &DetailsOptions{SleepAfter: 20}),
				&HttpRodResponsePipeline{
					BasePipeline: &forms.BasePipeline{
						Options: forms.GetNotRetribleOptions(),
						MetaInfo: &label.MetaInfo{
							Name: "http response from rod emulator",
						},
					},
					Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface,
						browser *rod.Browser, page *rod.Page, body *string, doc *goquery.Document) (error, context.Context) {

						logger.Warning(doc.Html())

						return nil, context
					},
					Cn: func(err error,
						context context.Context,
						task interfaces.TaskInterface,
						logger interfaces.LoggerInterface) {

						logger.Error("--- ", err.Error())
					},
				},
			},
		}

		newTask := fake.CreateDefaultTask()

		newTask.SetFetcherUrl("https://httpbin.org/ip")
		newTask.SetFetcherMethod("GET")
		newTask.Fetcher.Timeout = 60
		newTask.SetProxyServerUrl("http://localhost:8987")
		ctxGroup := Task.NewTaskCtx(newTask)

		err := g1.Run(ctxGroup)

		So(err, ShouldBeNil)
	})
}
