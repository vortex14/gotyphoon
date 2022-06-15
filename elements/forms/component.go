package forms

import (
	"path/filepath"

	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/folder"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Component struct {
	singleton.Singleton
	awaitabler.Object
	label.MetaInfo

	FileExt  string
	Language string
	Folder   *folder.Folder

	ProjectPath   string
	ProjectFolder *folder.Folder

	LOG interfaces.LoggerInterface

	Active      bool
	IsDebug     bool
	IsException bool

	Producers      interfaces.Producers
	Pipelines      interfaces.Pipelines
	Consumers      interfaces.Consumers
	QueuesSettings interfaces.Queue
}

func (c *Component) IsActive() bool {
	return c.Active
}

func (c *Component) Start(project interfaces.Project) {

}

func (c *Component) Close(project interfaces.Project) {

}

func (c *Component) AddPromise() {
	c.Add()
}

func (c *Component) PromiseDone() {
	c.Done()
}

func (c *Component) WaitPromises() {
	c.Await()
}

func (c *Component) InitFolder(componentName string) {

	c.Folder = &folder.Folder{Path: filepath.Join(c.ProjectPath, "project", componentName)}
	c.ProjectFolder = &folder.Folder{Path: c.ProjectPath}
}

func (c *Component) InitProducers() {

}

func (c *Component) StopConsumers() {

}

func (c *Component) StopProducers() {

}

func (c *Component) RunQueues() {

}

func (c *Component) InitConsumers(project interfaces.Project) {
	config := project.LoadConfig()
	queueSettings := config.TyComponents.Fetcher.Queues
	color.Yellow("current fetcher settings %+v", queueSettings)
	color.Yellow("InitConsumers for %s", c.Name)
}

func (c *Component) Restart(project interfaces.Project) {

}
