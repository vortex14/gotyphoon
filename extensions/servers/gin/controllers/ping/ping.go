package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GraphExtension "github.com/vortex14/gotyphoon/extensions/forms/graph"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	GinGraphExt "github.com/vortex14/gotyphoon/extensions/servers/gin/graph"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME 		= "ping"
	PATH        = "ping"
	DESCRIPTION = "Ping-Pong controller extension for Typhoon server"

	PING 		= "Ping"
	PONG 		= "Pong"
)

type PongResponse struct {
	Message string `json:"message"`
	Method string `json:"method"`
	Status bool	`json:"status"`
	Path string `json:"path"`
}

// handler
// @Tags main
// @Accept  json
// @Produce  json
// @Summary Ping Controller
// @Description Typhoon Ping Controller
// @Success 200 {object} PongResponse
// @Router /ping [get]
// @Router /ping [put]
// @Router /ping [post]
// @Router /ping [patch]
// @Router /ping [delete]
func handler (ctx *gin.Context, logger interfaces.LoggerInterface ) {
	logger.Debug(PING)
	ctx.JSON(200, gin.H{
		"method":  ctx.Request.Method,
		"message": PONG,
		"status":  true,
		"Path":    ctx.Request.RequestURI,
	})
}

var Controller = &GinExtension.Action{
	Action: &forms.Action{
		MetaInfo: &label.MetaInfo{
			Name: NAME,
			Path: PATH,
			Description: DESCRIPTION,
		},
		Methods: []string{interfaces.GET, interfaces.PATCH, interfaces.POST, interfaces.PUT, interfaces.DELETE},
	},
	GinController: handler,
}

var GraphController = &GinGraphExt.Action{
	Action: &GraphExtension.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name: NAME,
				Path: PATH,
				Description: DESCRIPTION,
			},
			Methods: []string{interfaces.GET, interfaces.PATCH, interfaces.POST, interfaces.PUT, interfaces.DELETE},
		},
	},
	GinController: handler,
}
