package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"

	"github.com/vortex14/gotyphoon/interfaces"

	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"

	"github.com/vortex14/gotyphoon/integrations/swagger"
)

type File struct {
	File []byte `json:"file" binding:"required"`
}

var FileController = &GinExtension.Action{
	Action: &forms.Action{
		BodyRequestModel: forms.BaseModelRequest{
			RequestModel: &File{}, Required: true,
			Type: swagger.FORMDATA,
		},
		MetaInfo: &label.MetaInfo{
			Tags:        []string{"File"},
			Path:        "file",
			Name:        "file",
			Description: "file server",
		},
		ResponseModels: map[int]interface{}{
			200: &TokenResponse{},
		},
		Methods: []string{interfaces.POST},
	},
	GinController: func(ctx *gin.Context, logger interfaces.LoggerInterface) {
		ctx.String(200, "asd")
	},
}
