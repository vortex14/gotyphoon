package home

import (
	"github.com/vortex14/gotyphoon/extensions/servers/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
)

type TyphoonMainResource struct {
	*interfaces.Resource
}

const (
	PING    = "ping"
	STATS   = "stats"
	METRICS = "metrics"
)

func Constructor() interfaces.ResourceInterface {
	return &TyphoonMainResource{
		Resource: &interfaces.Resource{
			Path: "/",
			Name: "Main resource of /",
			Description: `Main Resource of Typhoon Server. Access pass to / . 
                  This resource handle all action going after / . As example /ping, /stats, /metrics. 
                  Every action execute some controller which got ctx *gin.Context and logger as arguments.  
				  The Resource grouping all actions and all sub resources contain all other paths. `,
			Resource:    make(map[string]*interfaces.Resource),
			Middlewares: make([]*interfaces.Middleware, 0),
			Actions: map[string]*interfaces.Action{
				PING: ping.Controller,
			},
		},
	}
}
