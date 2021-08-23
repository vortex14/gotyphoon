package ghosts

import (
	"github.com/vortex14/gotyphoon/interfaces"
)

type DemonBuilder struct {
	Constructor func(opt *DemonOptions, project interfaces.Project) DemonInterface
	Options *DemonOptions
}

type DemonOptions struct {
	Required 	bool
	Name 		string
	Description string
	Sigil 		map[string]	interfaces.SigilBuilder
}

type DemonDecree struct {
	Name string
	Sigil string
	Cantrip string
}

type Demon struct {
	Options DemonOptions
	Project interfaces.Project
	Furies 	map[string]	FuryInterface
	Sigils 	map[string]	interfaces.SigilInterface
}

type DemonInterface interface {
	Run() error
	Sleep() error
	GetName() string
	AddSigil(sig *interfaces.Sigil) bool
	Execute (options *DemonDecree) error
}