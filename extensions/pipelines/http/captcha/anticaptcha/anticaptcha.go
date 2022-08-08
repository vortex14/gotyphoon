package anticaptcha

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"os"

	"github.com/nuveo/anticaptcha"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/interfaces"
)

func CreateAntiCaptchaPipeline() *pipelines.TaskPipeline {
	return &pipelines.TaskPipeline{
		BasePipeline: &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name:        "Anticaptcha",
				Description: "Decoding using the service anticaptcha",
			},
			Options: forms.GetCustomRetryOptions(2),
		},
		Fn: func(context context.Context,
			task interfaces.TaskInterface, logger interfaces.LoggerInterface) (error, context.Context) {

			b64cdata := GetData(context)

			token := task.GetSaveData(CaptchaKEY)

			if len(token) == 0 {
				return ApiKeyNotFound, context
			}

			c := &anticaptcha.Client{APIKey: os.Getenv(token)}
			logger.Info("send captcha")
			code, err := c.SendImage(
				b64cdata, // the image file encoded to base64
			)
			if err != nil {
				return err, context
			}

			task.SetSaveData(CaptchaCode, code)

			return nil, context
		},
	}
}
