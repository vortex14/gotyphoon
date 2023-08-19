package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	CheckRedisPath        = "check_redis"
	CheckRedisDescription = "Health check controller of redis"

	LocalhostRedisHost = "localhost"
	LocalhostRedisPort = 6379
)

var redisService = &redis.Service{
	Project: nil,
	Config: &interfaces.ServiceRedis{
		Name: "Test",
		Details: interfaces.RedisDetails{
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
func RedisHandler(ctx *gin.Context, logger interfaces.LoggerInterface) {
	ctx.JSON(200, &ServiceResponse{
		Status: redisService.Ping(),
	})
}

var RedisController = &GinExtension.Action{
	Action: &forms.Action{
		MetaInfo: &label.MetaInfo{
			Name:        CheckRedisPath,
			Description: CheckRedisDescription,
			Tags:        []string{"Services"},
		},
		Methods: []string{interfaces.GET},
	},
	GinController: RedisHandler,
}
