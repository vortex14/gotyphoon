package home

import (
	"github.com/vortex14/gotyphoon/extensions/servers/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

type TyphoonMainResource struct {
	*server.Resource
}

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

func Constructor(path string) server.ResourceInterface {
	return &TyphoonMainResource{
		Resource: &server.Resource{
			Path: path,
			Name: NAMEDefault,
			Description: DESCRIPTIONDefault,
			Resource:    make(map[string]*server.Resource),
			Middlewares: make([]*server.Middleware, 0),
			Actions: map[string]*server.Action{
				PING: ping.Controller,
			},
		},
	}
}