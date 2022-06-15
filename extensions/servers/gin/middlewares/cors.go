package middlewares

import (
	"context"
	Errors "github.com/vortex14/gotyphoon/errors"

	Gin "github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

func ConstructorCorsMiddleware(options server.CorsOptions) interfaces.MiddlewareInterface  {
	return &GinExtension.RequestMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Name:        "cors",
				Label:       "cors",
				Required:    true,
				Description: "cors middleware",
			},
		},
		Fn: func(
			context context.Context,
			ginCtx *Gin.Context,
			logger interfaces.LoggerInterface,
			reject func(err error),
			next func(ctx context.Context)) {

			logger.Info("check cors")

			ginCtx.Writer.Header().Set("Access-Control-Allow-Origin", options.AccessControlAllowOrigin)
			ginCtx.Writer.Header().Set("Access-Control-Allow-Headers", options.AccessControlAllowHeaders)
			ginCtx.Writer.Header().Set("Access-Control-Allow-Methods", options.AccessControlAllowMethods)
			ginCtx.Writer.Header().Set("Access-Control-Allow-Credentials", options.AccessControlAllowCredentials)

			if ginCtx.Request.Method == "OPTIONS" { ginCtx.AbortWithStatus(204); reject(Errors.ForceSkipRequest) }

		},
	}
}
