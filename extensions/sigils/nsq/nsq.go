package nsq

import (
	"github.com/vortex14/gotyphoon/errors"
	"os"

	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/extensions/sigils/nsq/cantrips"
	"github.com/vortex14/gotyphoon/interfaces"
)

const NAME = "*** NSQ SIGIL ***"
const DESCRIPTION = "NSQ Sigil"

type NSQ struct {
	Project interfaces.Project
	*interfaces.Sigil
}

var Options = &interfaces.SigilOptions{
	Name:        NAME,
	Required:    true,
	Description: DESCRIPTION,
	AllowedCantrips: []string{cantrips.PING},
	CantripsMap: map[string]interfaces.CantripInterface{
		cantrips.PING: cantrips.Create(cantrips.PING),
	},
}

func (n *NSQ) IsCritical() bool {
	return n.Critical
}

func (n *NSQ) Conjure(project interfaces.Project, name string) error {
	var err error
	switch name {
	case cantrips.PING:
		cantrip := n.Options.CantripsMap[name]
		err = cantrip.Conjure(project)
		if err != nil {
			return err
		}
		color.Green(">>> Cantrip %s done <<<", cantrip.GetName())
		return nil

	}
	return errors.NotFoundCantrip
}

func (n *NSQ) Crash(err error)  {
	color.Red("Crash cantrip: %s. Error: %s. Critical %t", n.Options.Name, err.Error(), n.IsCritical())
	if n.IsCritical() {
		os.Exit(1)
	}
}

func (n *NSQ) GetName() string {
	return n.Options.Name
}

func (n *NSQ) GetCantrips() map[string]interfaces.CantripInterface {
	return n.Options.CantripsMap
}

func Create(opt *interfaces.SigilOptions, project interfaces.Project) interfaces.SigilInterface {
	options := Options
	if opt != nil {
		options = opt
	}

	return &NSQ{
		Project: project,
		Sigil: &interfaces.Sigil{
			Options: options,
			Critical: opt.Required,
		},
	}
}

