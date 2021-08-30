package service

import (
	"sync"

	"github.com/fatih/color"
	NSQSig "github.com/vortex14/gotyphoon/extensions/sigils/nsq"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/ghosts"
)

var once sync.Once

const NAME = "--- Service Demon ----"
const DESCRIPTION = "Working with communication services"

type Demon struct {
	*ghosts.Demon
}

var demon *Demon

var BaseOptions = &ghosts.DemonOptions{
	Required: true,
	Name: NAME,
	Description: DESCRIPTION,
	Sigil: map[string]interfaces.SigilBuilder{
		interfaces.NSQ: {
			Constructor: NSQSig.Create,
			Options: NSQSig.Options,
		},
	},
}

func (d *Demon) GetName() string {
	return d.Options.Name
}

func (d *Demon) init() {
	if d.Options != nil && len(d.Options.Sigil) > 0 {
		d.Demon.Sigils = map[string]interfaces.SigilInterface{}
		for _, sigilBuilder := range d.Options.Sigil {
			sigil := sigilBuilder.Constructor(sigilBuilder.Options, d.Project)
			d.Sigils[sigil.GetName()] = sigil
			color.Yellow("|> %s initialized", sigil.GetName())
		}
	}
}



func (d *Demon) Run() error {
	d.init()
	return nil
}

func (d *Demon) Sleep() error {
	return nil
}


func (d *Demon) AddSigil(sig *interfaces.Sigil) bool {
	color.Red("%s", Errors.DemonFoolish.Error())
	return false
}

func Constructor(opt *ghosts.DemonOptions, project interfaces.Project) ghosts.DemonInterface {
	if opt == nil {
		opt = BaseOptions
	}

	once.Do(func() {
		demon = &Demon{
			Demon: &ghosts.Demon{
				Options: opt,
				Project: project,
			},
		}
	})

	return demon

}

func (d *Demon) Execute(detail *ghosts.DemonDecree) error {
	var err error
	if d.Project == nil {
		return Errors.DemonHasNotProject
	}

	if d.Sigils == nil || len(d.Sigils) == 0 {
		return Errors.DemonExecutingWithoutSettings
	}

	//color.Yellow("Execute >>> %+v", detail.Sigil)
	//color.Yellow("%+v", d.Sigils)
	if sig := d.Sigils[detail.Sigil]; sig != nil && d.Project != nil {
		err = sig.Conjure(d.Project, detail.Cantrip)
		if err != nil{
			return err
		}
	} else {
		err = Errors.NotFoundSigil
	}

	return err
}
