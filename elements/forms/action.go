package forms

import (
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	graphvizExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	"github.com/vortex14/gotyphoon/interfaces"

)

type Stats struct {
	input int64
}

type Action struct {
	*label.MetaInfo
	Stats

	Path           string
	Methods        [] string  // just yet HTTP Methods
	AllowedMethods [] string
	handlerPath    string
	graph          interfaces.GraphInterface
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

func (a *Action) SetHandlerPath(path string)  {
	a.handlerPath = path
}


func (a *Action) GetHandlerPath() string {
	return a.handlerPath
}

func (a *Action) UpdateGraphLabel(method string, path string)  {
	a.input ++
	println(a.input, method, path)
}

func (a *Action) AddMethodNodes()  {
	for _, method := range a.GetMethods() {
		a.graph.AddNode(&interfaces.NodeOptions{
			Name: graphvizExt.FormatBottomSpace(method),
			Shape: graphvizExt.SHAPEAction,
			EdgeOptions: &interfaces.EdgeOptions{
				NodeB:  a.handlerPath,
				ArrowS: 0.5,
			},
		})
	}
}

func (a *Action) SetGraph(parent interfaces.GraphInterface)  {
	a.graph = parent

	a.graph.AddNode(&interfaces.NodeOptions{
		Name: a.handlerPath,
		Shape: graphvizExt.SHAPETab,
	})

	a.AddMethodNodes()
}