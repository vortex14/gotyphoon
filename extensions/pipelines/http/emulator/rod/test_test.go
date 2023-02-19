package rod

import (
	"github.com/vortex14/gotyphoon/elements/forms"
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
		//_task.SetProxyAddress("http://154.53.89.38:8800")
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
