package forms

import (
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Action struct {
	*label.MetaInfo

	Path           string
	Methods        [] string  // just yet HTTP Methods
	AllowedMethods [] string
	Controller     interfaces.Controller // Controller of Action
	Pipeline       interfaces.PipelineGroupInterface
	PyController   interfaces.Controller // Python Controller Bridge of Action
	Middlewares    [] interfaces.MiddlewareInterface // Before a call to action we need to check this into middleware. May be client state isn't ready for serve
}

func (a *Action) AddMethod(name string) {
	logrus.Error(Errors.ActionAddMethodNotImplemented.Error())
}

func (a *Action) GetMiddlewareStack() [] interfaces.MiddlewareInterface {
	return a.Middlewares
}

func (a *Action) GetMethods() []string {
	return a.Methods
}

func (a *Action) GetController() interfaces.Controller {
	return a.Controller
}

func (a *Action) GetPipeline() interfaces.PipelineGroupInterface {
	return a.Pipeline
}