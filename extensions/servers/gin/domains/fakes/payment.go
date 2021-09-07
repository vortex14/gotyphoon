package fakes

import (
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

func CreatePaymentAction() interfaces.ActionInterface {
	return &GinExtension.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name:        NAME,
				Path:        FakePaymentPath,
				Description: "Fake Typhoon payment",
			},
			Methods :    []string{interfaces.GET},
		},
		GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
			ctx.JSON(200, fake.CreatePayment())
		},
	}
}
