package net_http

import (
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"

	"context"
	"errors"
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

		if len(sURL) == 0 {
			reject(UrlError)
			return
		}

		v, err := url.Parse(sURL[0])
		if err != nil || len(v.Host) == 0 {
			reject(UrlInvalidError)
			return
		}

	},
}
