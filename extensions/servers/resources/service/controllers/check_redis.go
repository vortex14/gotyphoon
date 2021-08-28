package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	CheckRedisPath = "check_redis"
	CheckRedisDescription = "Health check controller of redis"

	LocalhostRedisHost = "localhost"
	LocalhostRedisPort = 6379

)

var redisService = &redis.Service{
	Project: nil,
	Config:  &interfaces.ServiceRedis{
		Name: "Test",
		Details: struct {
			Host     string      `yaml:"host"`
			Port     int         `yaml:"port"`
			Password interface{} `yaml:"password"`
		}{
			Host: LocalhostRedisHost,
			Port: LocalhostRedisPort,
		},
	},
}

// RedisHandler
// @Tags Services
// @Produce  json
// @Summary controller of redis
// @Description Health check controller of redis
// @Success 200 {object} ServiceResponse
// @Router /api/v1/services/check_redis [get]
func RedisHandler (logger *logrus.Entry, ctx *gin.Context ) {
	ctx.JSON(200, &ServiceResponse{
		Status: redisService.Ping(),
	})
}

var RedisController = &interfaces.Action{
	Name: CheckRedisPath,
	Description: CheckRedisDescription,
	Controller: RedisHandler,
	Methods : []string{interfaces.GET},
}

