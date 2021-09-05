package component

import (
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/home"
	"github.com/vortex14/gotyphoon/interfaces"
)


func Constructor(name string, description string) interfaces.ResourceInterface {
	mainResource := home.Constructor("/").Get()
	mainResource.SetName(name)
	mainResource.SetDescription(description)
	return mainResource
}