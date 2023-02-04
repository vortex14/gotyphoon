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
	"time"
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
				CreateProxyRodRequestPipeline(
					forms.GetCustomRetryOptions(1, time.Duration(1)*time.Second),
					&DetailsOptions{SleepAfter: 1, MustElement: "#shitcoin > .message"},
				),
				&HttpRodResponsePipeline{
					BasePipeline: &forms.BasePipeline{
						NotIgnorePanic: true,
						Options:        forms.GetNotRetribleOptions(),
						MetaInfo: &label.MetaInfo{
							Name: "http response from rod emulator",
						},
					},
					Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface,
						browser *rod.Browser, page *rod.Page, body *string, doc *goquery.Document) (error, context.Context) {

						html, err := doc.Find("#shitcoin").Html()
						if err != nil {
							return err, context
						}
						logger.Debug(html)

						//logger.Warning(doc.Html())

						//defer page.MustClose()
						//defer browser.MustClose()

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

		newTask.SetFetcherUrl("https://honeypot.is/ethereum?address=0x879c61c813147627fe3ddb824f681f65550f2139")
		newTask.SetFetcherMethod("GET")
		newTask.Fetcher.Timeout = 60
		newTask.SetProxyServerUrl("http://localhost:8987")
		ctxGroup := Task.NewTaskCtx(newTask)

		err := g1.Run(ctxGroup)

		So(err, ShouldBeNil)
	})
}
