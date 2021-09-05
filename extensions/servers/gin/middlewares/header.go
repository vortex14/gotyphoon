package middlewares

import (
	"context"
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME = "Request Header Middleware"
	DESCRIPTION = "Setting request header for Typhoon task"
)

func ConstructorRequestHeaderMiddleware(required bool) interfaces.MiddlewareInterface {
	return &gin.RequestMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Name:        "check",
				Required:    true,
				Description: "check header",
			},
		},
		Fn: func(context context.Context, ginCtx *Gin.Context, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
			if len(ginCtx.Request.Header.Get("typhoon")) == 0 {
				reject(Errors.BadRequest)
			}
		},
	}
}