package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME        = "swagger"
	PATH        = "docs"
	DESCRIPTION = "Get swagger.json for Typhoon server"

	DOCS = "Docs"
)

func handler(ctx *gin.Context, server interfaces.ServerInterface, logger interfaces.LoggerInterface) {
	_, _ = ctx.Writer.Write(server.GetDocs())
}

var Controller = &GinExtension.Action{
	Action: &forms.Action{
		MetaInfo: &label.MetaInfo{
			Name:        NAME,
			Path:        PATH,
			Description: DESCRIPTION,
		},
		Methods: []string{http.MethodGet},
	},
	GinSController: handler,
}
