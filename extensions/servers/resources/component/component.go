package component

import (
	"github.com/vortex14/gotyphoon/extensions/servers/resources/home"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

type TyphoonComponentResource struct {
	*server.Resource
}

func Constructor(name string, description string) *TyphoonComponentResource {
	mainResource := home.Constructor().Get()
	mainResource.Name = name
	mainResource.Description = description
	return &TyphoonComponentResource{
		Resource: mainResource,
	}
}





