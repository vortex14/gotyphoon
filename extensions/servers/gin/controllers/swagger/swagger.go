package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
			Tags:        []string{"Docs"},
		},
		Methods: []string{http.MethodGet},
	},
	GinSController: handler,
}

var ControllerUI = &GinExtension.Action{
	Action: &forms.Action{
		MetaInfo: &label.MetaInfo{
			Name:        NAME,
			Path:        "swagger",
			Description: DESCRIPTION,
			Tags:        []string{"Docs"},
		},
		Methods: []string{http.MethodGet},
	},
	GinSController: handler,
}

func CreateUIController() {
	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:12735/docs"),
		ginSwagger.DefaultModelsExpandDepth(-1))

}
