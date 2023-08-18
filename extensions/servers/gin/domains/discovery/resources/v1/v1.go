package v1

import (
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/discovery/resources/v1/controllers"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/discovery/resources/v1/projects"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/home"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/service"
	"github.com/vortex14/gotyphoon/interfaces"
)

func Constructor() interfaces.ResourceInterface {
	return home.Constructor("/").
		AddAction(controllers.MeController).
		AddAction(controllers.LoginController).
		AddAction(controllers.FileController).
		AddResource(service.Constructor("services").Get()).
		AddResource(projects.Constructor("projects").Get())
}
