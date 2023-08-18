package projects

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/discovery/resources/v1/projects/controllers"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME        = "projects"
	DESCRIPTION = "resource of projects for discovery service"
)

func Constructor(path string) interfaces.ResourceInterface {
	return &forms.Resource{
		MetaInfo: &label.MetaInfo{
			Path:        path,
			Name:        NAME,
			Description: DESCRIPTION,
		},
		Actions: map[string]interfaces.ActionInterface{
			controllers.UnmoorControllerName:   controllers.UnmoorController,
			controllers.ProjectsControllerName: controllers.ProjectsController,
			controllers.RegisterControllerName: controllers.RegisterController,
		},
	}
}
