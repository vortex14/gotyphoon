package forms

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/integrations/swagger"
	"strconv"

	// /* ignore for building amd64-linux
	//	ghvzExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	// */

	"net/http"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type OnExit func()
type OnStart func(port int) error
type OnRequest func(context context.Context)
type OnServeHandler func(path string, method string, resource interfaces.ResourceInterface)
type OnResponse func(context context.Context, status int, data interfaces.Response)
type OnInitResource func(newResource interfaces.ResourceInterface)
type OnInitAction func(resource interfaces.ResourceInterface, action interfaces.ActionInterface)
type OnBuildSubResources func(subResource interfaces.ResourceInterface)
type OnBuildSubAction func(resource interfaces.ResourceInterface, action interfaces.ActionInterface)
type OnAddResource func(resource interfaces.ResourceInterface)
type OnReject func(status int, data interfaces.Response)
type OnCors func()

type ArchonChIN func(chan<- interface{})
type ArchonChOut func(<-chan interface{})

const (
	RoutePath      = "ROUTE_PATH"
	RequestContext = "REQUEST_CONTEXT"
)

type TyphoonServer struct {
	singleton.Singleton
	*label.MetaInfo

	Port      int
	Host      string
	Schema    string
	IsDebug   bool
	IsRunning bool

	ActiveSwagger bool

	Level string

	LOG         *logrus.Entry
	logInstance sync.Once
	Logger      *log.TyphoonLogger

	Resources map[string]interfaces.ResourceInterface

	OnStart             OnStart
	OnRequest           OnRequest
	OnInitAction        OnInitAction
	OnServeHandler      OnServeHandler
	OnBuildSubResources OnBuildSubResources
	OnBuildSubAction    OnBuildSubAction
	OnInitResource      OnInitResource
	OnAddResource       OnAddResource
	OnResponse          OnResponse
	OnReject            OnReject
	OnCors              OnCors
	OnExit              OnExit

	LoggerOptions  *log.Options
	TracingOptions *interfaces.TracingOptions

	ArchonChIN  ArchonChIN
	ArchonChOut ArchonChOut

	BuildGraph bool

	swagger *swagger.OpenApi

	// /* ignore for building amd64-linux
	//
	//	Graph           interfaces.GraphInterface
	//
	// */

}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func (s *TyphoonServer) GetDocs() []byte {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
	return nil
}

func (s *TyphoonServer) SetRouterGroup(resource interfaces.ResourceInterface, group interface{}) {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
}

func (s *TyphoonServer) Stop() error {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
	return Errors.ServerMethodNotImplemented
}

func (s *TyphoonServer) Restart() error {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
	return Errors.ServerMethodNotImplemented
}

func (s *TyphoonServer) RunServer(port int) error {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
	return Errors.ServerMethodNotImplemented
}

func (s *TyphoonServer) Init() interfaces.ServerInterface {
	s.LOG.Error("Init() ", Errors.ServerMethodNotImplemented.Error())
	return s
}

func (s *TyphoonServer) InitDocs() interfaces.ServerInterface {

	s.swagger = swagger.ConstructorNewFromArgs(
		s.Name,
		s.Description,
		s.Version,
		[]string{s.Schema, s.Host, strconv.Itoa(s.Port)})
	s.swagger.LOG = s.LOG
	return s
}

func (s *TyphoonServer) InitTracer() interfaces.ServerInterface {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
	return s
}

func (s *TyphoonServer) InitLogger() interfaces.ServerInterface {
	s.LOG = log.New(log.D{"server": s.Name})
	s.logInstance.Do(func() {
		if s.LoggerOptions == nil {
			return
		}
		s.Logger = (&log.TyphoonLogger{
			TracingOptions: s.TracingOptions,
			Name:           s.LoggerOptions.Name,
			Options:        *s.LoggerOptions,
		}).Init()
	})
	return s
}

//func (s *TyphoonServer) Init()  {
//
//}

func (s *TyphoonServer) Run() error {
	if !s.IsRunning && len(s.Resources) > 0 {
		if s.OnStart == nil {
			s.LOG.Error(Errors.ServerOnStartError.Error())
			return Errors.ServerOnStartError
		}
		if s.OnCors != nil {
			s.OnCors()
		}

		err := s.OnStart(s.Port)
		if err != nil {
			s.LOG.Error(fmt.Sprintf("Server %s, Error: %s", s.Name, err.Error()))
			return err
		}
		color.Yellow("Running Server %d ", s.Port)

	} else if len(s.Resources) == 0 {
		s.LOG.Error(Errors.NoResourcesAvailable.Error())
		return Errors.NoResourcesAvailable
	}

	return nil
}

func (s *TyphoonServer) InitResourcesMap() {
	s.Resources = make(map[string]interfaces.ResourceInterface)
}

func (s *TyphoonServer) InitRequestPath(context context.Context, path string) context.Context {
	return ctx.Update(context, RoutePath, path)
}

func (s *TyphoonServer) isMainAction(routePath string) bool {
	status := false

	paths := strings.Split(routePath, "/")

	if len(paths) == 2 {
		status = true
	}

	return status
}

// GetAction find the correct action for the client request
func (s *TyphoonServer) GetAction(
	requestPath string,
	logger interfaces.LoggerInterface,
	context context.Context,

) interfaces.ActionInterface {

	logger.Debug(fmt.Sprintf("get action for %s", requestPath))
	paths := strings.Split(requestPath, "/")

	var currentResource interfaces.ResourceInterface
	var currentAction interfaces.ActionInterface

	var joinedPath string
	var found bool

	for _, path := range paths {
		if s.isMainAction(requestPath) && currentResource == nil {
			currentResource = s.Resources["/"]
			found = true
			continue
		}

		if currentResource != nil {
			if ok, resource := currentResource.HasResource(path); ok {
				currentResource = resource
				found = true
				continue
			}
			if ok, action := currentResource.HasAction(path); ok {
				currentAction = action
				found = true
			}
		}

		if resource, ok := s.Resources[fmt.Sprintf("/%s", path)]; ok {
			currentResource = resource
			found = true
			continue
		}

		// For Main resource without home path on /
		if !found {
			if len(joinedPath) == 0 {
				joinedPath = path
			} else {
				joinedPath = fmt.Sprintf("%s/%s", joinedPath, path)
			}

			if resource, ok := s.Resources[fmt.Sprintf("/%s", joinedPath)]; ok {
				currentResource = resource
				found = true
				continue
			}
		}

	}
	return currentAction
}

// RunMiddlewareStack - pass request to controller middleware stack.
// Middleware may reject request by custom condition or just enrich context client request.
// Middleware may raise exception, but it be pass if flag required = false.
// Flag = true will be immediately reject client request
func (s *TyphoonServer) RunMiddlewareStack(
	requestContext context.Context,
	action interfaces.ActionInterface,
) (error, bool, context.Context) {

	statusMiddlewareStack := true
	var skipRequest bool
	var LastErrorMiddleware error
	{
		for _, middleware := range action.GetMiddlewareStack() {
			if !statusMiddlewareStack || skipRequest {
				break
			}
			middlewareLogger := log.Patch(s.LOG, log.D{"middleware": middleware.GetName()})
			//middlewareLogger := log.New(log.D{ "middleware": middleware.GetName()})

			// Refect client request
			middleware.Pass(requestContext, middlewareLogger, func(err error) {
				LastErrorMiddleware = err
				if middleware.IsRequired() {

					switch err {
					case Errors.ForceSkipRequest:
						skipRequest = true
					default:
						middlewareLogger.Error(err.Error())
						s.OnResponse(requestContext, http.StatusBadRequest, interfaces.Response{
							"message": err.Error(),
							"status":  false,
						})
						statusMiddlewareStack = false

					}

					return
				} else {
					middlewareLogger.Warning(err.Error())
				}
			}, func(context context.Context) {
				requestContext = context
			})

		}

	}

	return LastErrorMiddleware, statusMiddlewareStack, requestContext

}

func (s *TyphoonServer) initActions(resource interfaces.ResourceInterface) {
	for _, action := range resource.GetActions() {

		if len(action.GetMethods()) == 0 {
			s.LOG.Warning(Errors.ActionMethodsNotFound.Error())
			continue
		}

		for _, method := range action.GetMethods() {
			var handlerPath string
			if resource.GetPath() != "/" {
				handlerPath = fmt.Sprintf("%s/%s", resource.GetPath(), action.GetPath())
			} else {
				handlerPath = fmt.Sprintf("/%s", action.GetPath())
			}
			action.SetHandlerPath(handlerPath)

			actionLogger := log.Patch(s.LOG, log.D{"resource": resource.GetName(), "action": action.GetName()})
			action.SetLogger(actionLogger)
			s.LOG.Debug(fmt.Sprintf("need serve path: %s method: %s", handlerPath, method))

			s.AddSwaggerOperation(resource, action, method, handlerPath)

			s.initHandler(method, handlerPath, resource)
			if s.OnInitAction != nil {
				s.OnInitAction(resource, action)
			}
		}
	}
}

func (s *TyphoonServer) AddSwaggerResponses(action interfaces.ActionInterface, operation *openapi3.Operation) {

	errorResponseTitle := "Error response"
	s.swagger.AddSwaggerResponse(&errorResponseTitle, 422, action, operation, &ErrorResponse{})

	for status, model := range action.GetResponseModels() {
		responseTitle := "response"

		s.swagger.AddSwaggerResponse(&responseTitle, status, action, operation, model)
	}

}

func (s *TyphoonServer) AddSwaggerOperation(
	resource interfaces.ResourceInterface,
	action interfaces.ActionInterface,
	method, path string,
) {

	operation := s.swagger.AddSwaggerOperation(resource, action, method, path)

	s.AddSwaggerResponses(action, operation)

}

func (s *TyphoonServer) buildSubResources(parentPath string, newResource interfaces.ResourceInterface) {
	if newResource.GetCountSubResources() > 0 {
		for resourceName, subResource := range newResource.GetResources() {
			var resourcePath string
			if parentPath != "/" {
				resourcePath = fmt.Sprintf("%s/%s", parentPath, resourceName)
			} else {
				resourcePath = fmt.Sprintf("/%s", resourceName)
			}

			s.LOG.Debug("init subresource ", resourcePath, newResource.GetName(), newResource)

			subResource.SetLogger(s.LOG)

			subResource.SetPath(resourcePath)

			if s.OnBuildSubResources != nil {
				s.OnBuildSubResources(subResource)
			}

			s.buildSubActions(resourcePath, subResource)
			s.buildSubResources(resourcePath, subResource)

		}
	}
}

func (s *TyphoonServer) buildSubActions(parentPath string, newResource interfaces.ResourceInterface) {
	if newResource.GetCountActions() > 0 {
		for name, action := range newResource.GetActions() {
			for _, method := range action.GetMethods() {
				handlerPath := fmt.Sprintf("%s/%s", parentPath, name)
				s.LOG.Debug("init sub action ", handlerPath)
				s.AddSwaggerOperation(newResource, action, method, handlerPath)
				action.SetHandlerPath(handlerPath)
				if s.OnBuildSubAction != nil {
					s.OnBuildSubAction(newResource, action)
				}
				s.initHandler(method, handlerPath, newResource)
			}
		}
	}
}

func (s *TyphoonServer) initResource(newResource interfaces.ResourceInterface) error {
	if newResource.GetPath() == "" {
		return Errors.ResponsePathError
	}
	if s.Resources == nil {
		s.InitResourcesMap()
	}

	logger := log.Patch(s.LOG, log.D{"resource": newResource.GetName()})
	newResource.SetLogger(logger)

	if _, ok := s.Resources[newResource.GetPath()]; ok {
		return Errors.ResourceAlreadyExist
	} else {
		if s.OnInitResource != nil {
			s.OnInitResource(newResource)
		}
		s.Resources[newResource.GetPath()] = newResource
		s.initActions(newResource)

		// build resource fractal
		s.buildSubResources(newResource.GetPath(), newResource)
	}
	return nil
}

func (s *TyphoonServer) initHandler(method string, path string, resource interfaces.ResourceInterface) {
	if s.OnServeHandler == nil {
		s.LOG.Error(Errors.ServerOnHandlerMethodNotImplemented.Error())
	} else {
		s.OnServeHandler(method, path, resource)
	}
}

func (s *TyphoonServer) CreateResource(path string, opts label.MetaInfo) (error, interfaces.ResourceInterface) {
	newResource := &Resource{
		MetaInfo: &label.MetaInfo{
			Path:        path,
			Name:        opts.Name,
			Description: opts.Description,
		},
		Middlewares: make([]interfaces.MiddlewareInterface, 0),
		Actions:     make(map[string]interfaces.ActionInterface),
	}
	err := s.initResource(newResource)
	return err, newResource
}

func (s *TyphoonServer) AddResource(resource interfaces.ResourceInterface) interfaces.ServerInterface {
	logger := log.Patch(s.LOG, log.D{"resource": resource.GetName()})
	resource.SetLogger(logger)

	if s.OnAddResource != nil {
		s.OnAddResource(resource)
	}
	err := s.initResource(resource)
	if err != nil {
		color.Red("%s", err.Error())
	}
	return s
}

type ServerBuilder struct {
	Constructor func(project interfaces.Project) interfaces.ServerInterface
	server      interfaces.ServerInterface
	once        sync.Once
}

func (s *ServerBuilder) Run(project interfaces.Project) interfaces.ServerInterface {
	s.once.Do(func() { s.server = s.Constructor(project) })
	return s.server
}

// /* ignore for building amd64-linux
//
//func (s *TyphoonServer) InitGraph() interfaces.ServerInterface {
//	s.Graph = (&ghvzExt.Graph{
//		Options: &interfaces.GraphOptions{
//			IsCluster: true,
//		},
//		MetaInfo: &label.MetaInfo{
//			Name: fmt.Sprintf("Graph of %s",s.Name),
//			Label: s.Name,
//		},
//		Layout: ghvzExt.LAYOUTCirco,
//	}).Init()
//	return s
//}
//
//func (s *TyphoonServer) GetGraph() interfaces.GraphInterface {
//	return s.Graph
//}
//
//func (s *TyphoonServer) AddNewGraphResource(newResource interfaces.ResourceGraphInterface)  {
//	if s.Graph != nil {
//		subGraph := s.Graph.AddSubGraph(&interfaces.GraphOptions{
//			Name:      newResource.GetName(),
//			Label:     newResource.GetName(),
//			IsCluster: true,
//		})
//		s.LOG.Debug(fmt.Sprintf("init subGraph for %s", newResource.GetName()), subGraph)
//		newResource.SetGraph(subGraph)
//	} else {
//		s.LOG.Error("not found server graph. ",newResource.GetPath(), newResource.GetName())
//	}
//}
//
// */

func (s *TyphoonServer) GetServerEngine() interface{} {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
	return Errors.ServerMethodNotImplemented
}

func (s *TyphoonServer) SetServerEngine(server interface{}) {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error())
}

func (s *TyphoonServer) GetSwagger() []byte {
	return s.swagger.GetDump()

}
