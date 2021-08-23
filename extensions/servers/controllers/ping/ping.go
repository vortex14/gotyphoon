package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/interfaces"
)

func Handler (logger *logrus.Entry, ctx *gin.Context ) {
	logger.Debug("Ping")
	ctx.JSON(200, gin.H{
		"method": ctx.Request.Method,
		"message": "Pong",
		"status": true,
		"Path": ctx.Request.RequestURI,
	})
}

var Controller = &interfaces.Action{
	Methods : []string{interfaces.GET, interfaces.PATCH, interfaces.POST, interfaces.PUT, interfaces.DELETE},
	Name: "Ping Controller",
	Description: "Ping-Pong extension for Typhoon server",
	Controller: Handler,
}