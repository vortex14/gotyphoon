package rod

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"

	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"

	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func TestCreateDebugLauncher(t *testing.T) {

	browser := rod.New().ControlURL(CreateLauncher(Options{Debug: true}).MustLaunch())

	browser = browser.DefaultDevice(devices.IPadPro).Timeout(60 * time.Second).MustConnect()
	_ = browser.MustPage("https://www.wikipedia.org/")
	time.Sleep(5 * time.Second)
	_ = browser.Close()
}

func TestCreateBrowser(t *testing.T) {
	opts := Options{
		Timeout: 60 * time.Second,
		Debug:   true,
	}
	log.InitD()

	browser := CreateBaseBrowser(opts).MustConnect()

	_ = browser.MustPage("https://www.wikipedia.org/")
	time.Sleep(5 * time.Second)
	_ = browser.Close()

}

func TestCreateRodPipeline(t *testing.T) {

	Convey("create a new pipeline", t, func() {
		_task := fake.CreateDefaultTask()
		_task.SetProxyAddress("http://154.53.91.14:8800")
		_task.SetProxyServerUrl("http://proxy-manager.typhoon-s1.ru")
		p := CreateRodRequestPipeline(
			forms.GetNotRetribleOptions(),
			&DetailsOptions{
				ProxyRequired: true,
				Options: Options{
					Debug:   true,
					Device:  devices.IPadPro,
					Timeout: 600 * time.Second,
				},
				SleepAfter: 10,
			})

		_task.SetFetcherUrl("https://2ip.ru/")
		ctx := task.NewTaskCtx(_task)
		ctx = log.NewCtx(ctx, log.New(map[string]interface{}{"pipeline": "rod-request"}))

		var err error
		p.Run(ctx, func(pipeline interfaces.BasePipelineInterface, _err error) {
			err = _err
		}, func(ctx context.Context) {

		})

		So(err, ShouldBeNil)
	})

}

func TestRetryResponse(t *testing.T) {
	Convey("Move by coords", t, func() {
		count := 0
		g1 := forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{
				Name: "Rod group",
			},
			Stages: []interfaces.BasePipelineInterface{
				CreateRodRequestPipeline(
					forms.GetCustomRetryOptions(2, time.Duration(3)*time.Second),
					&DetailsOptions{
						SleepAfter: 0,
						Options: Options{
							Debug: true,
						},
					},
				),
				&HttpRodResponsePipeline{
					BasePipeline: &forms.BasePipeline{
						NotIgnorePanic: true,
						Options: &forms.Options{
							Retry: forms.RetryOptions{
								MaxCount: 5, Delay: time.Duration(3) * time.Second,
							},
						},
						MetaInfo: &label.MetaInfo{
							Name: "http response from rod emulator",
						},
					},
					Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface,
						browser *rod.Browser, page *rod.Page, body *string, doc *goquery.Document) (error, context.Context) {
						if count == 2 {
							return nil, context
						}
						count += 1
						return errors.New("a new error"), context
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

		_task := fake.CreateDefaultTask()
		_task.SetFetcherUrl("https://google.com/")
		_task.SetFetcherTimeout(600)
		ctxGroup := task.NewTaskCtx(_task)

		e := g1.Run(ctxGroup)

		So(count, ShouldEqual, 2)

		So(e, ShouldBeNil)
	})
}
