package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

const (
	NAME 		= "ping"
	DESCRIPTION = "Ping-Pong controller extension for Typhoon server"

	PING 		= "Ping"
	PONG 		= "Pong"
)

type PongResponse struct {
	Method string `json:"method"`
	Message string `json:"message"`
	Status bool	`json:"status"`
	Path string `json:"path"`
}

// handler
// @Tags main
// @Accept  json
// @Produce  json
// @Summary Ping Controller
// @Description Typhoon Ping Controller
// @Success 200 {object} PongResponse
// @Router /ping [get]
// @Router /ping [put]
// @Router /ping [post]
// @Router /ping [patch]
// @Router /ping [delete]
func handler (logger *logrus.Entry, ctx *gin.Context ) {
	logger.Debug(PING)
	ctx.JSON(200, gin.H{
		"method": ctx.Request.Method,
		"message": PONG,
		"status": true,
		"Path": ctx.Request.RequestURI,
	})
}

var Controller = &server.Action{
	Name: NAME,
	Description: DESCRIPTION,
	Controller: handler,
	Methods : []string{server.GET, server.PATCH, server.POST, server.PUT, server.DELETE},
}
