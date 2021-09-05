package ghosts

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/interfaces"
)

type DemonCallOptions struct {
	Demon string
	Sigil string
	Cantrip string
}

type ArchonInterface interface {
	AddPromise()
	ClosePromise()
	AwaitDecision()
	Call(detail *DemonDecree) error
	RunDemons(project interfaces.Project)
	RunProjectServers(project interfaces.Project)
}

type Archon struct {
	Name 		string
	Description string
	Demons 		[] *DemonBuilder
	Observers 	map[string]	ObserverInterface
	Sisyphuses 	map[string]	SisyphusInterface
	Servers 	map[string]	*forms.ServerBuilder
}

type ProtoArchon struct {
	Archons []ArchonInterface
}

