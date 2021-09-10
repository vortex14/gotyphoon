package forms

import (
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/utils"
	"net/http"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	ghvzExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type OnExit           func()
type OnStart          func(port int) error
type OnRequest        func(context context.Context)
type OnResponse       func(status int, data interfaces.Response)
type OnReject         func(status int, data interfaces.Response)
type OnServeHandler   func(path string, method string)

type ArchonChIN       func(chan<- interface{} )
type ArchonChOut      func(<-chan interface{} )

const (
	RoutePath      = "ROUTE_PATH"
	RequestContext = "REQUEST_CONTEXT"
)

type TyphoonServer struct {
	*label.MetaInfo

	Port            int
	IsDebug         bool
	IsRunning   	bool
	BuildGraph      bool
	Level           string
	Instance        sync.Once

	LOG             *logrus.Entry
	logInstance     sync.Once
	Logger          *log.TyphoonLogger

	Resources   	map[string]interfaces.ResourceInterface

	OnStart         OnStart
	OnRequest       OnRequest
	OnServeHandler  OnServeHandler
	OnResponse      OnResponse
	OnReject        OnReject
	OnExit          OnExit

	LoggerOptions	*log.Options
	SwaggerOptions  *interfaces.SwaggerOptions
	TracingOptions  *interfaces.TracingOptions

	ArchonChIN      ArchonChIN
	ArchonChOut     ArchonChOut

	Graph           interfaces.GraphInterface

}

func (s *TyphoonServer) Stop() error  {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error()); return Errors.ServerMethodNotImplemented
}

func (s *TyphoonServer) Restart() error {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error()); return Errors.ServerMethodNotImplemented
}

func (s *TyphoonServer) RunServer(port int) error {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error()); return Errors.ServerMethodNotImplemented
}

func (s *TyphoonServer) Init() interfaces.ServerInterface {
	s.LOG.Error("Init() ",Errors.ServerMethodNotImplemented.Error()); return s
}

func (s *TyphoonServer) InitGraph() interfaces.ServerInterface {
	s.Graph = (&ghvzExt.Graph{
		BaseGraph: &ghvzExt.BaseGraph{
			Layout: ghvzExt.LAYOUTCirco,
			MetaInfo:  &label.MetaInfo{
				Name: fmt.Sprintf("Graph of %s",s.Name),
			},
		},
	}).Init()
	return s
}

func (s *TyphoonServer) GetGraph() interfaces.GraphInterface {
	return s.Graph
}

func (s *TyphoonServer) InitDocs() interfaces.ServerInterface {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error()); return s
}

func (s *TyphoonServer) InitTracer() interfaces.ServerInterface {
	s.LOG.Error(Errors.ServerMethodNotImplemented.Error()); return s
}

func (s *TyphoonServer) InitLogger() interfaces.ServerInterface {
	if s.IsDebug { log.InitD() }
	s.LOG = log.New(log.D{"server": s.Name})
	s.logInstance.Do(func() {
		if s.LoggerOptions == nil { return }
		s.Logger = (&log.TyphoonLogger{
			TracingOptions: s.TracingOptions,
			Name: s.LoggerOptions.Name,
			Options: log.Options{ BaseOptions: &log.BaseOptions{} },
		}).Init()
	})
	return s
}

func (s *TyphoonServer) Run() error {
	if !s.IsRunning && len(s.Resources) > 0 {
		if s.OnStart == nil { s.LOG.Error(Errors.ServerOnStartError.Error()); return Errors.ServerOnStartError }
		err := s.OnStart(s.Port)
		if err != nil {
			s.LOG.Error(fmt.Sprintf("Server %s, Error: %s", s.Name, err.Error()))
			return err
		}
		color.Yellow("Running Server %d ", s.Port)

	} else if len(s.Resources) == 0 {
		s.LOG.Error(Errors.NoResourcesAvailable.Error()); return Errors.NoResourcesAvailable
	}

	return nil
}

func (s *TyphoonServer) InitResourcesMap()  {
	s.Resources = make(map[string]interfaces.ResourceInterface)
}

func (s *TyphoonServer) InitRequestPath(context context.Context, path string) context.Context  {
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
	) (error, bool, context.Context)  {

	statusMiddlewareStack := true
	var LastErrorMiddleware error
	{
		for _, middleware := range action.GetMiddlewareStack() {
			if !statusMiddlewareStack { break }
			middlewareLogger := log.New(log.D{ "middleware": middleware.GetName()})

			// Refect client request
			middleware.Pass(requestContext, middlewareLogger, func(err error) {
				LastErrorMiddleware = err
				if middleware.IsRequired() {
					middlewareLogger.Error(err.Error())
					s.OnResponse(http.StatusBadRequest, interfaces.Response{
						"message": err.Error(),
						"status": false,
					})
					statusMiddlewareStack = false
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

func (s *TyphoonServer) initActions(resource interfaces.ResourceInterface)  {
	for _, action := range resource.GetActions() {

		if len(action.GetMethods()) == 0 { s.LOG.Warning(Errors.ActionMethodsNotFound.Error()); continue }

		for _, method := range action.GetMethods() {
			var handlerPath string
			if resource.GetPath() != "/" {
				handlerPath = fmt.Sprintf("%s/%s",resource.GetPath(), action.GetPath())
			} else {
				handlerPath = fmt.Sprintf("/%s", action.GetPath())
			}
			action.SetHandlerPath(handlerPath)
			s.LOG.Debug(fmt.Sprintf("need serve path: %s", handlerPath))

			s.initHandler(method, handlerPath)

			if utils.NotNill(s.Graph) { resource.AddGraphActionNode(action) }
		}
	}
}

func (s *TyphoonServer) buildSubResources(parentPath string, newResource interfaces.ResourceInterface)  {
	for resourceName, subResource := range newResource.GetResources() {
		resourcePath := fmt.Sprintf("%s/%s", parentPath, resourceName)

		if subResource.GetCountActions() > 0 {
			s.buildSubActions(resourcePath, subResource)
		}

		if subResource.GetCountSubResources() > 0 {
			s.buildSubResources(resourcePath, subResource)
		}
	}
}

func (s *TyphoonServer) buildSubActions(parentPath string, newResource interfaces.ResourceInterface)  {
	for name, action := range newResource.GetActions() {
		for _, method := range action.GetMethods() {
			handlerPath := fmt.Sprintf("%s/%s", parentPath, name)
			s.initHandler(method, handlerPath)
		}
	}
}

func (s *TyphoonServer) initResource(newResource interfaces.ResourceInterface) error {
	if newResource.GetPath() == "" {
		return Errors.ResponsePathError
	}
	if s.Resources == nil { s.InitResourcesMap() }

	logger := log.Patch(s.LOG, log.D{"resource": newResource.GetName()})
	newResource.SetLogger(logger)

	if _, ok := s.Resources[newResource.GetPath()]; ok {
		return Errors.ResourceAlreadyExist
	} else {
		if s.Graph != nil {
			s.LOG.Debug(fmt.Sprintf("init subGraph for %s", newResource.GetName()))
			subGraph := s.Graph.AddSubGraph(&interfaces.GraphOptions{
				Name:            newResource.GetName(),
				Label:           newResource.GetName(),
				IsCluster:       true,
			})
			newResource.SetGraph(subGraph)
		}
		s.Resources[newResource.GetPath()] = newResource
		s.initActions(newResource)

		// build resource fractal
		if newResource.GetCountSubResources() > 0 {
			s.buildSubResources(newResource.GetPath(), newResource)
		}
	}
	return nil
}

func (s *TyphoonServer) initHandler(method string, path string)  {
	if s.OnServeHandler == nil { s.LOG.Error(Errors.ServerOnHandlerMethodNotImplemented.Error()) } else {
		s.OnServeHandler(method, path)
	}
}

func (s *TyphoonServer) CreateResource(path string, opts label.MetaInfo) (error, interfaces.ResourceInterface) {
	newResource := &Resource{
		MetaInfo: &label.MetaInfo{
			Path:            path,
			Name:            opts.Name,
			Description:     opts.Description,
		},
		Middlewares:     make([]interfaces.MiddlewareInterface, 0),
		Actions:         make(map[string]interfaces.ActionInterface),
	}
	err := s.initResource(newResource)
	return err, newResource
}

func (s *TyphoonServer) AddResource(resource interfaces.ResourceInterface) interfaces.ServerInterface {
	if s.BuildGraph && s.Graph == nil { s.InitGraph() }
	logger := log.Patch(s.LOG, log.D{"resource": resource.GetName()})
	resource.SetLogger(logger)
	err := s.initResource(resource)
	if err != nil { color.Red("%s", err.Error()) }
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