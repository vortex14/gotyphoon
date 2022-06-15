package test

import (
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME = "test"
	DESCRIPTION = "Typhoon test resource"
)

func Constructor() interfaces.ResourceInterface {
	return &forms.Resource{
		MetaInfo: &label.MetaInfo{
			Path:        NAME,
			Name:        NAME,
			Description: DESCRIPTION,
		},
		Resources:   make(map[string]interfaces.ResourceInterface),
		Middlewares: make([]interfaces.MiddlewareInterface, 0),
		Actions: map[string]interfaces.ActionInterface{
			NAME: ping.Controller,

			"demo": &GinExtension.Action{
				Action: &forms.Action{
					MetaInfo: &label.MetaInfo{
						Name:        "demo",
						Path:        "demo",
						Description: "demo",
					},
					Methods: []string{interfaces.GET, interfaces.PATCH, interfaces.POST, interfaces.PUT, interfaces.DELETE},
				},
				GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
					ctx.JSON(200, Gin.H{"status": true, "msg": "demo"})
				},
			},
		},
	}
}