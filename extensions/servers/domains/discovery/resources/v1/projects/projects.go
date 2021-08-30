package projects

import (
	"github.com/vortex14/gotyphoon/extensions/servers/domains/discovery/resources/v1/projects/controllers"
	"github.com/vortex14/gotyphoon/interfaces"

	"github.com/vortex14/gotyphoon/elements/forms"
)

const (
	NAME = "projects"
	DESCRIPTION = "resource of projects for discovery service"
)

func Constructor(path string) interfaces.ResourceInterface {
	return &forms.Resource{
		Path: path,
		Name: NAME,
		Description: DESCRIPTION,
		Resources:    make(map[string]interfaces.ResourceInterface),
		Middlewares: make([]interfaces.MiddlewareInterface, 0),
		Actions: map[string]interfaces.ActionInterface{
			controllers.ProjectsControllerName: controllers.ProjectsController,
			controllers.RegisterControllerName: controllers.RegisterController,
			controllers.UnmoorControllerName: controllers.UnmoorController,
		},
	}
}
