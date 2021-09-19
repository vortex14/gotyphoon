package graph

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME 		= "graph"
	PATH        = "server-graph"
	DESCRIPTION = "Graphviz controller extension for visualization Typhoon server"

	GRAPH 		= "Graph"
)

// handler
// @Tags main
// @Router /graph [get]
func handler (ctx *gin.Context, server interfaces.ServerInterface, logger interfaces.LoggerInterface ) {
	// /* ignore for building amd64-linux
	graphS := server.(interfaces.ServerGraphInterface)
	serverGraph := graphS.GetGraph()

	exportFormat := ctx.Request.URL.Query().Get("format")

	switch exportFormat {
	case "dot", "jpg", "svg": _, _ = ctx.Writer.Write(serverGraph.Render(exportFormat))
	case "all": ctx.JSON(200, gin.H{"formats": [...]string{"dot", "jpg", "svg"}})
	default   : ctx.JSON(400, gin.H{"error": Errors.GraphNotFoundFormat.Error()})
	}

	// */
}

var Controller = &GinExtension.Action{
	Action: &forms.Action{
		MetaInfo: &label.MetaInfo{
			Name: NAME,
			Path: PATH,
			Description: DESCRIPTION,
		},
		Methods: []string{interfaces.GET},
	},
	GinSController: handler,
}

