package v1

import (
	"github.com/vortex14/gotyphoon/extensions/servers/domains/discovery/resources/v1/controllers"
	"github.com/vortex14/gotyphoon/extensions/servers/domains/discovery/resources/v1/projects"
	"github.com/vortex14/gotyphoon/extensions/servers/resources/home"
	"github.com/vortex14/gotyphoon/extensions/servers/resources/service"
	"github.com/vortex14/gotyphoon/interfaces"
)

func Constructor() interfaces.ResourceInterface {
	return home.Constructor("/api/v1").
		AddAction(controllers.MeController).
		AddAction(controllers.LoginController).
		AddResource(service.Constructor("services").Get()).
		AddResource(projects.Constructor("projects").Get())
}