package home

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/graph"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/swagger"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/test"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAMEDefault        = "Main resource of /"
	DESCRIPTIONDefault = `Main Resource of Typhoon Server. Access pass to / . 
                  This resource handle all action going after / . As example /ping, /stats, /metrics. 
                  Every action execute some controller which got ctx *gin.Context and logger as arguments.  
				  The Resource grouping all actions and all sub resources contain all other paths. `

	PING    = "ping"
	STATS   = "stats"
	METRICS = "metrics"
	GRAPH   = "server-graph"
	DOCS    = "docs"
	SWAGGER = "swagger"
)

func Constructor(path string) interfaces.ResourceInterface {
	return (&forms.Resource{
		MetaInfo: &label.MetaInfo{
			Path:        path,
			Name:        NAMEDefault,
			Description: DESCRIPTIONDefault,
		},
		Actions: map[string]interfaces.ActionInterface{
			PING:  ping.Controller,
			GRAPH: graph.Controller,
			DOCS:  swagger.Controller,
			//SWAGGER: swagger.ControllerUI,
		},
	}).AddResource(test.Constructor())
}
