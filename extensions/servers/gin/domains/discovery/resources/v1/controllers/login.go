package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"

	"github.com/vortex14/gotyphoon/interfaces"

	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
)

const (
	NAME        = "login"
	DESCRIPTION = "Discovery login controller extension for Typhoon server"

	JWTAUTHDefault = "eyJsb2dpbiI6InR5cGhvb24iLCJlbWFpbCI6InR5cGhvb25AdHlwaG9vbi1zMS5ydSIsInJvbGVzIjpbXSwiYWxnIjoiSFMy+" +
		"NTYifQ.e30.7m63q7oIzRooWceOw5DX-S8av4NHx_AbQx8oibISgZU"
)

type UserPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// handler
// @Tags Auth
// @Accept  json
// @Produce  plain
// @Summary Discovery login controller
// @Description Typhoon Discovery login controller
// @Success 200 {string}
// @Router /api/v1/login [post]
func handler(ctx *gin.Context, logger interfaces.LoggerInterface) {
	ctx.String(200, JWTAUTHDefault)
}

var LoginController = &GinExtension.Action{
	Action: &forms.Action{
		BodyRequestModel: forms.BaseModelRequest{RequestModel: &UserPayload{}, Required: true},
		MetaInfo: &label.MetaInfo{
			Tags:        []string{"Auth"},
			Path:        NAME,
			Name:        NAME,
			Description: DESCRIPTION,
		},
		ResponseModels: map[int]interface{}{
			200: &TokenResponse{},
		},
		Methods: []string{interfaces.POST},
	},
	GinController: handler,
}
