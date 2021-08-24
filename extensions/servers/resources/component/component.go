package component

import (
	"github.com/vortex14/gotyphoon/extensions/servers/resources/home"
	"github.com/vortex14/gotyphoon/interfaces"
)

type TyphoonComponentResource struct {
	*interfaces.Resource
}

func Constructor(name string, description string) *TyphoonComponentResource {
	mainResource := home.Constructor().Get()
	mainResource.Name = name
	mainResource.Description = description
	return &TyphoonComponentResource{
		Resource: mainResource,
	}
}





