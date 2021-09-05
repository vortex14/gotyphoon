package fakes

import (
	Gin "github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

func CreateCategoryAction() interfaces.ActionInterface {
	return &GinExtension.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name:        NAME,
				Path:        FakeCategoryPath,
				Description: "Fake Typhoon category",
			},
			Methods :    []string{interfaces.GET},
		},
		GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
			ctx.JSON(200, fake.CreateCategory())
		},

	}
}

func CreateCategoriesAction() interfaces.ActionInterface {
	return &GinExtension.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name:        NAME,
				Path:        FakeCategoriesPath,
				Description: "Fake Typhoon categories",
			},
			Methods :    []string{interfaces.GET},
		},
		GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
			ctx.JSON(200, fake.CreateCategories())
		},

	}
}