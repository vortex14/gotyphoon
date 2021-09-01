package server

import (
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/forms"
	"net/http"
	"strings"
	"sync"
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/Bose/go-gin-opentracing"
	"github.com/fatih/color"
	Gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/logger"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"

	"github.com/vortex14/gotyphoon/extensions/servers/pipelines/gin"
)


type ServerBuilder struct {
	Constructor func(project interfaces.Project) interfaces.ServerInterface
	server      interfaces.ServerInterface
	once        sync.Once
}

func (s *ServerBuilder) Run(project interfaces.Project) interfaces.ServerInterface {
	s.once.Do(func() {
		s.server = s.Constructor(project)
	})
	return s.server
}

type TyphoonServer struct {
	Port 			int
	isRunning   	bool
	IsDebug         bool
	Level 			string
	server 			*Gin.Engine
	logger 			*logger.TyphoonLogger
	LOG             *logrus.Entry
	resources   	map [string]interfaces.ResourceInterface
	callbacks 		map [string]func(ctx *Gin.Context)

	TracingOptions  *interfaces.TracingOptions
	LoggerOptions	*interfaces.BaseLoggerOptions
	SwaggerOptions  *interfaces.SwaggerOptions


	*interfaces.BaseServerLabel

}

func (s *TyphoonServer) InitLogger() interfaces.ServerInterface {
	if s.IsDebug { log.InitD() }
	s.LOG = log.New(log.D{"server": s.Name})
	if s.LoggerOptions != nil {
		s.logger = &logger.TyphoonLogger{
			TracingOptions: s.TracingOptions,
			Name: s.LoggerOptions.Name,
			Options: logger.Options{
				BaseLoggerOptions: s.LoggerOptions,
			},
		}
		s.logger.Init()
	}
	return s
}

func (s *TyphoonServer) InitDocs() interfaces.ServerInterface {
	if s.SwaggerOptions != nil {
		url := ginSwagger.URL(s.SwaggerOptions.DocEndpoint)
		s.LOG.Info("InitDocs URL >>> %s", s.SwaggerOptions.DocEndpoint)
		s.server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	return s

}

func (s *TyphoonServer) InitTracer() interfaces.ServerInterface {
	if s.logger != nil && s.TracingOptions != nil {
		p := ginopentracing.OpenTracer([]byte(s.logger.GetTracerHeader()))

		s.server.Use(p)

		useBanner := s.TracingOptions.UseBanner
		useUTC :=  s.TracingOptions.UseUTC

		s.server.Use(ginlogrus.WithTracing(logrus.StandardLogger(),
			useBanner,
			time.RFC3339,
			useUTC,
			"requestID",
			[]byte("typhoon-trace-id"), // where jaeger might have put the trace id
			[]byte("RequestID"),     // where the trace ID might already be populated in the headers
			ginlogrus.WithAggregateLogging(false)))
	}

	return s
}

func (s *TyphoonServer) Init() interfaces.ServerInterface {
	if s.server == nil {
		s.resources = make(map[string]interfaces.ResourceInterface)
		s.callbacks = make(map[string]func(ctx *Gin.Context))

		s.server = Gin.New()
		s.server.Use(Gin.Recovery())
	}

	return s
}

func (s *TyphoonServer) Run() error {
	if !s.isRunning && len(s.resources) > 0 {
		port := fmt.Sprintf(":%d", s.Port)
		err := s.server.Run(port)
		if err != nil {
			s.LOG.Error("Server %s, Error: %s", s.Name, err.Error())
			return err
		}

		color.Yellow("Running Server %s ", port)

	} else if len(s.resources) == 0 {
		logrus.Error(Errors.NoResourcesAvailable.Error())
	}


	return nil
}

func (s *TyphoonServer) Stop() error  {


	return nil
}

func (s *TyphoonServer) Restart() error {
	return nil
}

func (s *TyphoonServer) isMainAction(ctx *Gin.Context) bool {
	status := false

	paths := strings.Split(ctx.Request.URL.Path, "/")

	if len(paths) == 2 {
		status = true
	}

	return status
}

// getAction find the correct action for the client request
func (s *TyphoonServer) getAction(logger interfaces.LoggerInterface, ctx *Gin.Context) (interfaces.ActionInterface) {
	actionPath := ctx.Request.URL.Path
	//logger.Debug(actionPath, s.resources)
	paths := strings.Split(actionPath, "/")

	var currentResource interfaces.ResourceInterface
	var currentAction interfaces.ActionInterface

	var joinedPath string
	var found bool


	for _, path := range paths {
		if s.isMainAction(ctx) && currentResource == nil {
			currentResource = s.resources["/"]
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

		if resource, ok := s.resources[fmt.Sprintf("/%s", path)]; ok {
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

			if resource, ok := s.resources[fmt.Sprintf("/%s", joinedPath)]; ok {
				currentResource = resource
				found = true
				continue
			}

		}

	}
	return currentAction
}

// requestHandler handle all HTTP request in here
func (s *TyphoonServer) requestHandler(ginContext *Gin.Context)  {
	logger := ginlogrus.GetCtxLogger(ginContext)
	action := s.getAction(logger, ginContext)
	if action == nil {
		logger.Error("Not found")
		ginContext.JSON(404, Gin.H{
			"message": "Not Found",
			"status": false,
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"controller": action.GetName(),
	})


	// Pass request to controller middleware stack.
	// Middleware may reject request by custom condition or just enrich context client request.
	// Middleware may raise exception, but it be pass if flag required = false.
	// Flag = true will be immediately reject client request
	statusMiddlewareStack := true
	var LastErrorMiddleware error
	{
		for _, middleware := range action.GetMiddlewareStack() {
			if !statusMiddlewareStack { break }
			middlewareLogger := log.New(log.D{ "middleware": middleware.GetName()})

			// Refect client request
			middleware.Pass(ginContext, middlewareLogger, func(err error) {
				LastErrorMiddleware = err
				if middleware.IsRequired() {
					middlewareLogger.Error(err.Error())

					ginContext.JSON(http.StatusBadRequest, Gin.H{
						"message": err.Error(),
						"status": false,
					})
					statusMiddlewareStack = false
					return
				} else {
					middlewareLogger.Warning(err.Error())
				}
			}, func(context context.Context) {

			})

		}

	}
	if statusMiddlewareStack {
		controller := action.GetController()
		if controller != nil {
			controller(ginContext, logger)
		} else if pipeline := action.GetPipeline(); pipeline != nil {
			mainCtx := ctx.Update(ctx.New(), gin.CTX, ginContext)
			pipeline.Run(mainCtx)
		}

	} else {
		color.Red("%s", LastErrorMiddleware.Error())
	}


}

func (s *TyphoonServer) initActions(resource interfaces.ResourceInterface)  {
	for _, action := range resource.GetActions() {

		if len(action.GetMethods()) == 0 { s.LOG.Warning(Errors.ActionMethodsNotFound.Error()); break }

		for _, method := range action.GetMethods() {
			var handlerPath string
			if resource.GetPath() != "/" {
				handlerPath = fmt.Sprintf("%s/%s",resource.GetPath(), action.GetPath())
			} else {
				handlerPath = fmt.Sprintf("/%s", action.GetPath())
			}
			logrus.Debug(fmt.Sprintf("serve path: %s", handlerPath))
			s.Serve(method, handlerPath, s.requestHandler)

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
			s.Serve(method, handlerPath, s.requestHandler)
		}
	}
}

func (s *TyphoonServer) initResource(newResource interfaces.ResourceInterface) error {
	if newResource.GetPath() == "" {
		return Errors.ResponsePathError
	}

	if _, ok := s.resources[newResource.GetPath()]; ok {
		return Errors.ResourceAlreadyExist
	} else {
		s.resources[newResource.GetPath()] = newResource
		s.initActions(newResource)

		// build resource fractal
		if newResource.GetCountSubResources() > 0 {
			s.buildSubResources(newResource.GetPath(), newResource)
		}
	}
	return nil
}

func (s *TyphoonServer) initHandler(method string, path string, callback func(ctx *Gin.Context))  {

	var handler func(ctx *Gin.Context)

	if callback == nil {
		handler = s.requestHandler
	} else {
		handler = callback
	}

	switch method {
	case interfaces.GET:
		s.server.GET(path, handler)
	case interfaces.POST:
		s.server.POST(path, handler)
	case interfaces.PUT:
		s.server.PUT(path, handler)
	case interfaces.PATCH:
		s.server.PATCH(path, handler)
	case interfaces.DELETE:
		s.server.DELETE(path, handler)
	}
}

func (s *TyphoonServer) Serve(method string, path string, callback func(ctx *Gin.Context))  {
	s.Init()
	if len(s.resources) == 0 {
		s.callbacks[path] = callback
		s.initHandler(method, path, callback)
	} else {
		s.initHandler(method, path, nil)
	}

}

func (s *TyphoonServer) CreateResource(path string, opts interfaces.BaseServerLabel) (error, interfaces.ResourceInterface) {
	newResource := &forms.Resource{
		Path:            path,
		Name:            opts.Name,
		Description:     opts.Description,
		Middlewares:     make([]interfaces.MiddlewareInterface, 0),
		Actions:         make(map[string]interfaces.ActionInterface),
	}
	err := s.initResource(newResource)
	return err, newResource
}

func (s *TyphoonServer) AddResource(resource interfaces.ResourceInterface) interfaces.ServerInterface {
	s.Init()
	err := s.initResource(resource)
	if err != nil {
		color.Red("%s", err.Error())
	}
	return s
}