package rucaptcha

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"time"

	api2captcha "github.com/vortex14/2captcha-go"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/interfaces"
)

func CreateRuCaptchaPipeline() *pipelines.TaskPipeline {
	return &pipelines.TaskPipeline{
		BasePipeline: &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name:        "Rucaptcha",
				Description: "Decoding using the service rucaptcha",
			},
			Options: forms.GetCustomRetryOptions(2, time.Duration(1)*time.Second),
		},
		Fn: func(context context.Context,
			task interfaces.TaskInterface, logger interfaces.LoggerInterface) (error, context.Context) {

			b64cdata := GetData(context)

			ruCap := api2captcha.Normal{
				Base64: b64cdata,
			}

			ruCap.CaseSensitive = true

			token := task.GetSaveData(CaptchaKEY)
			if len(token) == 0 {
				return ApiKeyNotFound, context
			}

			client := api2captcha.NewClient(task.GetSaveData(CaptchaKEY))

			code, err := client.Solve(ruCap.ToRequest())

			if err != nil {
				return err, context
			}

			task.SetSaveData(CaptchaCode, code)

			return nil, context
		},
	}
}
