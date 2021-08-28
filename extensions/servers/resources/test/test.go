package test

import (
	"github.com/vortex14/gotyphoon/extensions/servers/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

const (
	NAME = "test"
	DESCRIPTION = "Typhoon test resource"
)

type TyphoonTestResource struct {
	*server.Resource
}

func Constructor() interfaces.ResourceInterface {
	return &TyphoonTestResource{
		Resource: &server.Resource{
			Path: NAME,
			Name: NAME,
			Description: DESCRIPTION,
			Resource:    make(map[string]*server.Resource),
			Middlewares: make([]*server.Middleware, 0),
			Actions: map[string]*interfaces.Action{
				NAME: ping.Controller,
			},
		},
	}
}
