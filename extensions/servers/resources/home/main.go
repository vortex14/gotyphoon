package home

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/servers/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAMEDefault = "Main resource of /"
	DESCRIPTIONDefault = `Main Resource of Typhoon Server. Access pass to / . 
                  This resource handle all action going after / . As example /ping, /stats, /metrics. 
                  Every action execute some controller which got ctx *gin.Context and logger as arguments.  
				  The Resource grouping all actions and all sub resources contain all other paths. `

	PING    = "ping"
	STATS   = "stats"
	METRICS = "metrics"
)

func Constructor(path string) interfaces.ResourceInterface {
	return &forms.Resource{
		Path: path,
		Name: NAMEDefault,
		Description: DESCRIPTIONDefault,
		Resources:    make(map[string]interfaces.ResourceInterface),
		Middlewares: make([]interfaces.MiddlewareInterface, 0),
		Actions: map[string]interfaces.ActionInterface{
			PING: ping.Controller,
		},
	}
}