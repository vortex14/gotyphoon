package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"

	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
)

const (

	PATH                    = "me"
	NAMEMeController 		= "me"
	DESCRIPTIONMeController = "Me discovery JWT controller extension for Typhoon server"

	JWTAUTHLOGINDefault     = "typhoon"
	JWTAUTHEmailDefault     = "typhoon@typhoon-s1.ru"
	JWTRoleDefault          = "admin"
)

type MeResponse struct {
	Login string     `json:"login"`
	Email string    `json:"email"`
	Roles []string	`json:"roles"`
}

// handler
// @Tags Auth
// @Accept  json
// @Produce  plain
// @Summary Discovery me controller
// @Description Typhoon Discovery me controller
// @Success 200 {object} MeResponse
// @Router /api/v1/me [get]
func meHandler (ctx *gin.Context, logger interfaces.LoggerInterface ) {
	ctx.JSON(200, &MeResponse{
		Login: JWTAUTHLOGINDefault,
		Email: JWTAUTHEmailDefault,
		Roles: []string{JWTRoleDefault},
	})
}

var MeController = &GinExtension.Action{
	Action: &forms.Action{
		MetaInfo: &label.MetaInfo{
			Path:        PATH,
			Name:        NAMEMeController,
			Description: DESCRIPTIONMeController,
		},
		Methods :    []string{interfaces.GET},
	},
	GinController:  meHandler,
}