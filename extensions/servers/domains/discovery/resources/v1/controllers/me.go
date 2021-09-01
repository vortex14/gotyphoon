package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/interfaces"
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

var MeController = &interfaces.Action{
	Path: PATH,
	Name: NAMEMeController,
	Description: DESCRIPTIONMeController,
	Controller: meHandler,
	Methods : []string{interfaces.GET},
}