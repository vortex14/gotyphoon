package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/interfaces/server"
)


const (
	NAME 		= "login"
	DESCRIPTION = "Discovery login controller extension for Typhoon server"

	JWTAUTHDefault = "eyJsb2dpbiI6InR5cGhvb24iLCJlbWFpbCI6InR5cGhvb25AdHlwaG9vbi1zMS5ydSIsInJvbGVzIjpbXSwiYWxnIjoiSFMy+" +
		"NTYifQ.e30.7m63q7oIzRooWceOw5DX-S8av4NHx_AbQx8oibISgZU"

)


// handler
// @Tags Auth
// @Accept  json
// @Produce  plain
// @Summary Discovery login controller
// @Description Typhoon Discovery login controller
// @Success 200 {string}
// @Router /api/v1/login [post]
func handler (logger *logrus.Entry, ctx *gin.Context ) {
	ctx.String(200, JWTAUTHDefault)
}

var LoginController = &server.Action{
	Name: NAME,
	Description: DESCRIPTION,
	Controller: handler,
	Methods : []string{server.POST},
}

