package test

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/servers/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME = "test"
	DESCRIPTION = "Typhoon test resource"
)

func Constructor() interfaces.ResourceInterface {
	return &forms.Resource{
		Path: NAME,
		Name: NAME,
		Description: DESCRIPTION,
		Resources:    make(map[string]interfaces.ResourceInterface),
		Middlewares: make([]interfaces.MiddlewareInterface, 0),
		Actions: map[string]interfaces.ActionInterface{
			NAME: ping.Controller,
		},
	}
}
