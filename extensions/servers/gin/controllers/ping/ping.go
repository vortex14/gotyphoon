package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"net/http"

	// /* ignore for building amd64-linux
	//
	//	GraphExtension "github.com/vortex14/gotyphoon/extensions/forms/graph"
	//	GinGraphExt "github.com/vortex14/gotyphoon/extensions/servers/gin/graph"
	//
	// */

	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME        = "ping"
	PATH        = "ping"
	DESCRIPTION = "Ping-Pong controller extension for Typhoon server"

	PING = "Ping"
	PONG = "Pong"
)

type PongResponse struct {
	Message string `json:"message"`
	Method  string `json:"method"`
	Status  bool   `json:"status"`
	Path    string `json:"path"`
}

func handler(ctx *gin.Context, logger interfaces.LoggerInterface) {
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
		ResponseModels: map[int]interface{}{200: &PongResponse{}},
		MetaInfo: &label.MetaInfo{
			Name:        NAME,
			Path:        PATH,
			Description: DESCRIPTION,
			Tags:        []string{"System"},
		},
		Methods: []string{http.MethodGet},
	},
	GinController: handler,
}

// /* ignore for building amd64-linux
//
//var GraphController = &GinGraphExt.Action{
//	Action: &GraphExtension.Action{
//		Action: &forms.Action{
//			MetaInfo: &label.MetaInfo{
//				Name: NAME,
//				Path: PATH,
//				Description: DESCRIPTION,
//			},
//			Methods: []string{interfaces.GET, interfaces.PATCH, interfaces.POST, interfaces.PUT, interfaces.DELETE},
//		},
//	},
//	GinController: handler,
//}
//
// */
