package forms

import (
	"context"
	"github.com/sirupsen/logrus"

	// /* ignore for building amd64-linux
//	"fmt"
//	graphvizExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	// */
	"github.com/vortex14/gotyphoon/log"

	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Stats struct {
	Input int64
}

type Action struct {
	*label.MetaInfo
	LOG            interfaces.LoggerInterface
	Stats

	Path           string
	Methods        [] string   //just yet HTTP Methods
	AllowedMethods [] string
	handlerPath    string

	Controller     interfaces.Controller  //Controller of Action
	Pipeline       interfaces.PipelineGroupInterface
	PyController   interfaces.Controller  //Python Controller Bridge of Action
	Middlewares    [] interfaces.MiddlewareInterface  //Before a call to action we need to check this into middleware. May be client state isn't ready for serve

	// /* ignore for building amd64-linux
//	Graph          interfaces.GraphInterface
	// */

}

func (a *Action) AddMethod(name string) {
	logrus.Error(Errors.ActionAddMethodNotImplemented.Error())
}

func (a *Action) IsPipeline() bool {
	status := true
	if a.Pipeline == nil {
		status = false
	}
	return status
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



func (a *Action) InitPipelineGraph()  {
	pipelineLogger := log.Patch(a.LOG.(*logrus.Entry), log.D{"pipeline-group": a.GetPipeline().GetName()})
	a.Pipeline.SetLogger(pipelineLogger)
	// /* ignore for building amd64-linux
//	a.Pipeline.SetGraph(a.Graph)
//	a.Pipeline.InitGraph(a.GetHandlerPath())
	// */
}

func (a *Action) UpdateGraphLabel(method string, path string)  {

	///* ignore for building amd64-linux
//	a.Input ++
//
//		labelAction := fmt.Sprintf(`
//	
//	R: %d
//	
//	`, a.Input)
//
//	if a.Graph != nil {
//		a.Graph.UpdateEdge(&interfaces.EdgeOptions{
//			NodeA: method,
//			NodeB: path,
//			LabelH: labelAction,
//			Color: graphvizExt.COLORNavy,
//
//		})
//	}
//
    //  */

}

func (a *Action) AddMethodNodes() {

	// /* ignore for building amd64-linux
//	for _, method := range a.GetMethods() {
//		a.Graph.AddNode(&interfaces.NodeOptions{
//			Name: method,
//			Label: graphvizExt.FormatBottomSpace(method),
//			Shape: graphvizExt.SHAPEAction,
//			EdgeOptions: &interfaces.EdgeOptions{
//				NodeB:  a.GetHandlerPath(),
//				ArrowS: 0.5,
//			},
//		})
//	}
//
	//*/
}


// /* ignore for building amd64-linux
//
//func (a *Action) SetGraph(parent interfaces.GraphInterface, buildMethods bool)  {
//	a.Graph = parent
//
//	opts := &interfaces.NodeOptions{
//		Name: a.GetHandlerPath(),
//		Label: graphvizExt.FormatSpace(a.GetHandlerPath()),
//		Shape: graphvizExt.SHAPETab,
//		Style: graphvizExt.StyleFilled,
//		BackgroundColor: graphvizExt.COLORGray,
//		EdgeOptions: &interfaces.EdgeOptions{},
//
//	}
//
//	a.Graph.AddNode(opts)
//
//	if buildMethods {
//		a.AddMethodNodes()
//	}
//
//	if a.IsPipeline() {
//		a.InitPipelineGraph()
//	}
//}
//
//func (a *Action) SetGraphNodes(nodes map[string]interfaces.NodeInterface)  {
//	println(fmt.Sprintf("ACTION NAME :%s INIT GRAPH NODES ------ >>> %+v", a.GetPath(), nodes))
//	a.Graph.SetNodes(nodes)
//
//	println(fmt.Sprintf("ACTION NAME :%s GET GRAPH NODES ------ >>> %+v", a.GetPath(), a.Graph.GetNodes()))
//}
//
//func (a *Action) GetGraph() interfaces.GraphInterface {
//	return a.Graph
//}
//
//
// */

func (a *Action) OnRequest(method string, path string)  {
	// /* ignore for building amd64-linux
//	a.UpdateGraphLabel(method, path)
	//*/
}

func (a *Action) Run(ctx context.Context, logger interfaces.LoggerInterface)  {

}

func (a *Action) SetLogger(logger interfaces.LoggerInterface)  {
	a.LOG = logger
}