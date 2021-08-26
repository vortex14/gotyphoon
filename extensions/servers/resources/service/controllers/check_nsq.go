package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/integrations/nsq"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

const (
	CheckNSQPath = "check_nsq"
	CheckNSQDescription = "Health check controller of NSQ"

	LocalhostNSQDHost = "localhost"
	LocalhostNSQDPort = 4161

)


var nsqService = &nsq.Service{
	Config: &interfaces.ConfigProject{
		NsqlookupdIP: fmt.Sprintf("%s:%d", LocalhostNSQDHost, LocalhostNSQDPort),
	},
}

// NSQHandler
// @Tags Services
// @Produce  json
// @Summary controller of NSQ
// @Description Health check controller of NSQ
// @Success 200 {object} ServiceResponse
// @Router /api/v1/services/check_nsq [get]
func NSQHandler (logger *logrus.Entry, ctx *gin.Context ) {
	ctx.JSON(200, &ServiceResponse{
		Status: nsqService.Ping(),
	})
}

var NSQController = &server.Action{
	Name: CheckNSQPath,
	Description: CheckNSQDescription,
	Controller: NSQHandler,
	Methods : []string{server.GET},
}

