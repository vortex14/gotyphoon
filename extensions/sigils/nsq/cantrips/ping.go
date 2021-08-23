package cantrips

import (
	"github.com/fatih/color"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

const DESCRIPTION =  "Cantrip NSQ Ping"

type PingCantrip struct {
	Name string
	Description string
}

func (p *PingCantrip) Conjure(project interfaces.Project) error {
	project.LoadServices(interfaces.TyphoonIntegrationsOptions{NSQ: interfaces.MessageBrokerOptions{Active: true}})
	service := project.GetService(interfaces.NSQ)
	if !service.Ping() {
		color.Red("No ping to NSQ")
		return Errors.ServiceNotAvailable
	}
	return nil
}

func (p *PingCantrip) GetName() string {
	return p.Name
}

func Create(name string) interfaces.CantripInterface {
	switch name {
	case PING:
		return &PingCantrip{
			Name: PING,
			Description: DESCRIPTION,
		}
	}

	color.Red("%s", Errors.NotFoundCantrip)

	return nil

}