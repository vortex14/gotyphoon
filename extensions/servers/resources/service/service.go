package service

import (
	"github.com/vortex14/gotyphoon/extensions/servers/resources/service/controllers"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

const (
	NAME = "services"
	DESCRIPTION = "Main Service Resource"
)

type TyphoonServiceResource struct {
	*server.Resource
}

func Constructor(path string) server.ResourceInterface {
	return &TyphoonServiceResource{
		Resource: &server.Resource{
			Path: path,
			Name: NAME,
			Description: DESCRIPTION,
			Resource:    make(map[string]*server.Resource),
			Middlewares: make([]*server.Middleware, 0),
			Actions: map[string]*server.Action{
				controllers.CheckMongoPath: controllers.MongoController,
				controllers.CheckNSQPath: controllers.NSQController,
				controllers.CheckRedisPath: controllers.RedisController,
			},
		},
	}
}