package service

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/service/controllers"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME        = "services"
	DESCRIPTION = "Main Service Resource"
)

func Constructor(path string) interfaces.ResourceInterface {
	return &forms.Resource{
		MetaInfo: &label.MetaInfo{
			Path:        path,
			Name:        NAME,
			Description: DESCRIPTION,
		},
		Resources:   make(map[string]interfaces.ResourceInterface),
		Middlewares: make([]interfaces.MiddlewareInterface, 0),
		Actions: map[string]interfaces.ActionInterface{
			controllers.CheckMongoPath: controllers.MongoController,
			controllers.CheckNSQPath:   controllers.NSQController,
			controllers.CheckRedisPath: controllers.RedisController,
		},
	}
}
