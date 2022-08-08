package rucaptcha

import (
	"github.com/vortex14/gotyphoon/ctx"

	"context"
)

const (
	CaptchaBASE64 = "captcha_base64_data"
	CaptchaKEY    = "RUCAPTCHA_KEY"
	CaptchaCode   = "rucaptcha_code"
)

func PatchCtx(context context.Context, data string) context.Context {
	return ctx.Update(context, CaptchaBASE64, data)
}

func GetData(context context.Context) string {
	return ctx.Get(context, CaptchaBASE64).(string)
}
