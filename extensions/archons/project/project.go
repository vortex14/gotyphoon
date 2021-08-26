package project

import (
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces/server"
	"sync"

	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/ghosts"
)

type Archon struct {
	promise  sync.WaitGroup
	*ghosts.Archon
	demons map[string] ghosts.DemonInterface
}

func (a *Archon) AddPromise()  {
	a.promise.Add(1)
}

func (a *Archon) Call(detail *ghosts.DemonDecree) error {
	if demon := a.demons[detail.Name]; demon != nil {
		err := demon.Execute(detail)
		return err
	}
	return Errors.DemonNotFound
}

func (a *Archon) ClosePromise()  {
	a.promise.Done()
}

func (a *Archon) AwaitDecision()  {
	a.promise.Wait()
}


func (a *Archon) RunDemons(project interfaces.Project)  {
	a.demons = map[string]ghosts.DemonInterface{}
	for _, demonBuilder := range a.Demons {
		demon := demonBuilder.Constructor(demonBuilder.Options, project)
		if err := demon.Run(); err != nil {
			color.Red("%s", err.Error())
			continue
		}

		a.demons[demon.GetName()] = demon
		color.Yellow("%s awake", demon.GetName())
	}
}

func (a *Archon) RunProjectServers(project interfaces.Project)  {
	for _, it := range a.Servers {
		go func(server server.Interface) {
			if err:= server.Run(); err != nil {
				color.Red("%s", err.Error())
				return
			}
			server.InitDocs()
		}(it.Run(project))
	}
}