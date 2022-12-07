package net_http

import (
	"context"
	b64 "encoding/base64"
	"errors"
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
	"net/url"
)

var (
	UrlError        = errors.New("not found url as query param")
	UrlInvalidError = errors.New("not valid url")
)

var UrlRequiredMiddleware = &gin.RequestMiddleware{
	Middleware: &forms.Middleware{
		MetaInfo: &label.MetaInfo{
			Name:        "check",
			Required:    true,
			Description: "check request",
		},
	},
	Fn: func(context context.Context, ctx *Gin.Context, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
		params := ctx.Request.URL.Query()

		sURL := params["url"]
		encodeFlag := params["encode"]
		isBase64 := false

		if len(encodeFlag) > 0 && encodeFlag[0] == "base64" {
			isBase64 = true
		}

		if len(sURL) == 0 {
			logger.Error(params.Encode())
			reject(UrlError)
			return
		}
		URL := sURL[0]
		if isBase64 {

			sDec, be := b64.StdEncoding.DecodeString(URL)
			if be != nil {
				logger.Error(URL)
				reject(UrlError)
				return
			}
			URL = string(sDec)
			logger.Warning(URL)
		}

		v, err := url.Parse(URL)
		if err != nil || len(v.Host) == 0 {
			logger.Error(URL)
			reject(UrlInvalidError)
			return
		}

	},
}
