package forms

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/log"
	"go.uber.org/zap"

	"errors"

	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Stats struct {
	Input int64
}

type BaseModelRequest struct {
	RequestModel interface{}
	Required     bool
	Type         string
}

type HeaderRequestModel struct {
	ErrorModel      interface{}
	Model           interface{}
	ErrorStatusCode int
}

type Action struct {
	*label.MetaInfo
	LOG interfaces.LoggerInterface
	Stats

	Path           string
	Methods        []string //just yet HTTP Methods
	AllowedMethods []string
	handlerPath    string

	Service          interface{}
	Params           interface{}
	Headers          HeaderRequestModel
	BodyRequestModel BaseModelRequest
	ResponseModels   map[int]interface{}

	//Cn func(ctx context.Context, err error)

	Controller   interfaces.Controller //Controller of Action
	Pipeline     interfaces.PipelineGroupInterface
	PyController interfaces.Controller            //Python Controller Bridge of Action
	Middlewares  []interfaces.MiddlewareInterface //Before a call to action we need to check this into middleware. May be client state isn't ready for serve

	Cn func(
		err error,
		context context.Context,
		logger interfaces.LoggerInterface,
	)

	// /* ignore for building amd64-linux
	//	Graph          interfaces.GraphInterface
	// */

}

func (a *Action) AddMethod(name string) {
	logrus.Error(Errors.ActionAddMethodNotImplemented.Error())
}

func (a *Action) Cancel(ctx context.Context, logger interfaces.LoggerInterface, err error) {
	logger.Warn("Action.cancel", zap.Error(Errors.ActionFailed))
}

func (a *Action) IsPipeline() bool {
	status := true
	if a.Pipeline == nil {
		status = false
	}
	return status
}

func (a *Action) GetMiddlewareStack() []interfaces.MiddlewareInterface {
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

func (a *Action) SetHandlerPath(path string) {
	a.handlerPath = path
}

func (a *Action) GetHandlerPath() string {
	return a.handlerPath
}

func (a *Action) InitPipelineGraph() {
	pipelineLogger := log.Patch(a.LOG.(*zap.Logger), log.D{"pipeline-group": a.GetPipeline().GetName()})
	a.Pipeline.SetLogger(pipelineLogger)
	// /* ignore for building amd64-linux
	//	a.Pipeline.SetGraph(a.Graph)
	//	a.Pipeline.InitGraph(a.GetHandlerPath())
	// */
}

func (a *Action) UpdateGraphLabel(method string, path string) {

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

func (a *Action) OnRequest(method string, path string) {
	// /* ignore for building amd64-linux
	//	a.UpdateGraphLabel(method, path)
	//*/
}

func (a *Action) SafeRun(run func() error, catch func(err error)) {
	defer func() {

		if r := recover(); r != nil {
			panicE := errors.New(fmt.Sprintf("%s: %s", PanicException, r))
			catch(panicE)
		}

	}()

	if err := run(); err != nil {
		catch(err)
	}

}

func (a *Action) Run(ctx context.Context, logger interfaces.LoggerInterface) {

}

func (a *Action) SetLogger(logger interfaces.LoggerInterface) {
	a.LOG = logger
}

func (a *Action) GetRequestModel() interface{} {
	return a.BodyRequestModel.RequestModel
}

func (a *Action) GetHeadersModel() interface{} {
	return a.Headers.Model
}

func (a *Action) GetHeadersErrModel() interface{} {
	return a.Headers.ErrorModel
}

func (a *Action) GetErrorHeadersStatusCode() int {
	return a.Headers.ErrorStatusCode
}

func (a *Action) GetRequestType() string {
	return a.BodyRequestModel.Type
}

func (a *Action) IsRequiredRequestModel() bool {
	return a.BodyRequestModel.Required
}

func (a *Action) IsValidRequestBody() bool {
	status := false
	if a.GetRequestModel() != nil && a.IsRequiredRequestModel() {
		status = true
	}
	return status
}

func (a *Action) GetResponseModels() map[int]interface{} {
	return a.ResponseModels
}

func (a *Action) GetParams() interface{} {
	return a.Params
}

func (a *Action) GetService() interface{} {
	return a.Service
}
