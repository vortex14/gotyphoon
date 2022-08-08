package rucaptcha

import (
	b64 "encoding/base64"
	api2captcha "github.com/vortex14/2captcha-go"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"
	"os"
	"testing"

	"context"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	log.InitD()
}

func TestDecode(t *testing.T) {
	data := utils.ReadFile("data2.jpg")
	b64data := b64.StdEncoding.EncodeToString([]byte(data))

	Convey("check api key", t, func() {
		So(len(os.Getenv("RUCAPTCHA_KEY")) > 0, ShouldBeTrue)
	})

	Convey("decode", t, func() {
		client := api2captcha.NewClient(os.Getenv(CaptchaKEY))

		cap2 := api2captcha.Normal{
			Base64: b64data,
		}

		cap2.CaseSensitive = true
		println(client)
		//code, err := client.Solve(cap2.ToRequest())
		//So(err, ShouldBeNil)
		//fmt.Println("code " + code)
	})

}

func TestCaptchaPipeline(t *testing.T) {

	Convey("rucaptcha decode pipeline", t, func(c C) {
		newTask := fake.CreateDefaultTask()
		newTask.SetSaveData("RUCAPTCHA_KEY", os.Getenv(CaptchaKEY))

		ctx := task.NewTaskCtx(newTask)
		ctx = log.NewCtxValues(ctx, log.D{"log": true})

		data := utils.ReadFile("data2.jpg")
		b64data := b64.StdEncoding.EncodeToString([]byte(data))

		ctx = PatchCtx(ctx, b64data)

		p := CreateRuCaptchaPipeline()
		var E error
		var FinalCtx context.Context

		p.Run(ctx, func(pipeline interfaces.BasePipelineInterface, err error) {

			E = err

		}, func(ctx context.Context) {
			FinalCtx = ctx
		})

		So(E, ShouldBeNil)

		status, resultTask := task.Get(FinalCtx)

		So(status, ShouldBeTrue)

		code := resultTask.GetSaveData(CaptchaCode)

		So(code, ShouldEqual, "66W0Q")

	})

}
