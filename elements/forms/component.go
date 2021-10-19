package forms

import (
	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Component struct {
	singleton.Singleton
	awaitabler.Object
	label.MetaInfo

	FileExt string
	Language string


	Active bool
	isDebug bool
	isException bool

	Producers interfaces.Producers
	Pipelines interfaces.Pipelines
	Consumers interfaces.Consumers
	QueuesSettings interfaces.Queue
}

func (c *Component) Start()  {

}

func (c *Component) Stop()  {

}
