package fakes

import (
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	GraphExtension "github.com/vortex14/gotyphoon/extensions/forms/graph"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	GinGraphExt "github.com/vortex14/gotyphoon/extensions/servers/gin/graph"
	"github.com/vortex14/gotyphoon/interfaces"
)

func handler(ctx *Gin.Context, logger interfaces.LoggerInterface)  {
	ctx.JSON(200, fake.CreateUpc())
}

// /* ignore for building amd64-linux

var GraphController = &GinGraphExt.Action{
	Action: &GraphExtension.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name: NAME,
				Path: FakeUPCPath,
				Description: "Fake Typhoon UPC code",
			},
			Methods: []string{interfaces.GET},
		},
	},
	GinController: handler,

}

// */

func CreateUpcAction() interfaces.ActionInterface {
	return &GinExtension.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name:        NAME,
				Path:        FakeUPCPath,
				Description: "Fake Typhoon UPC code",
			},
			Methods :    []string{interfaces.GET},
		},
		GinController: handler,

	}
}

