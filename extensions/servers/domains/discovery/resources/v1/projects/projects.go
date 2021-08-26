package projects

import (
	"github.com/vortex14/gotyphoon/extensions/servers/domains/discovery/resources/v1/projects/controllers"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

const (
	NAME = "projects"
	DESCRIPTION = "resource of projects for discovery service"
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
				controllers.ProjectsControllerName: controllers.ProjectsController,
				controllers.RegisterControllerName: controllers.RegisterController,
				controllers.UnmoorControllerName: controllers.UnmoorController,
			},
		},
	}
}
