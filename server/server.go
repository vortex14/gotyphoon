package server

import (
	"fmt"
	"strings"
	"sync"
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/Bose/go-gin-opentracing"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/logger"
	"github.com/vortex14/gotyphoon/interfaces"
)


type ServerBuilder struct {
	Constructor func(project interfaces.Project) interfaces.ServerInterface
	server interfaces.ServerInterface
	once sync.Once
}

func (s *ServerBuilder) Run(project interfaces.Project) interfaces.ServerInterface {
	s.once.Do(func() {
		s.server = s.Constructor(project)
	})
	return s.server
}

type TyphoonServer struct {
	Port 			int
	IsRunning   	bool
	Level 			string
	server 			*gin.Engine
	logger 			*logger.TyphoonLogger
	resources   	map[string]*interfaces.Resource
	callbacks 		map [string]func(ctx *gin.Context)

	TracingOptions  *interfaces.TracingOptions
	LoggerOptions	*interfaces.BaseLoggerOptions
	SwaggerOptions *interfaces.SwaggerOptions


	*interfaces.BaseServerLabel

}

func (s *TyphoonServer) InitLogger() interfaces.ServerInterface {
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
		color.Red("InitDocs URL >>> %s", s.SwaggerOptions.DocEndpoint)
		s.server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	return s

}

func (s *TyphoonServer) InitTracer() interfaces.ServerInterface  {
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
		s.resources = make(map[string]*interfaces.Resource)
		s.callbacks = make(map[string]func(ctx *gin.Context))

		s.server = gin.New()
		s.server.Use(gin.Recovery())
	}

	return s
}

func (s *TyphoonServer) Run() error {
	if !s.IsRunning {
		port := fmt.Sprintf(":%d", s.Port)
		err := s.server.Run(port)
		if err != nil {
			color.Red("Server %s, Error: %s", s.Name, err.Error())
			return err
		}

		color.Yellow("Running Server %s ", port)

	}


	return nil
}

func (s *TyphoonServer) Stop() error  {


	return nil
}

func (s *TyphoonServer) Restart() error {
	return nil
}

func (s *TyphoonServer) isMainAction(ctx *gin.Context) bool {
	status := false

	paths := strings.Split(ctx.Request.URL.Path, "/")

	if len(paths) == 2 {
		status = true
	}

	return status
}

func (s *TyphoonServer) getAction(ctx *gin.Context) (*interfaces.Action) {
	actionPath := ctx.Request.URL.Path
	paths := strings.Split(actionPath, "/")
	//color.Yellow("%+v, %d", paths, len(paths))
	var currentResource *interfaces.Resource
	var currentAction *interfaces.Action
	for _, path := range paths {
		if s.isMainAction(ctx) && currentResource == nil {
			currentResource = s.resources["/"]
			continue
		}

		if currentResource != nil {
			if subResource, ok := currentResource.Resource[path]; ok {
				currentResource = subResource
				continue
			}
			if currentHandler, ok := currentResource.Actions[path]; ok {
				currentAction = currentHandler
			}
		}

		if resource, ok := s.resources[fmt.Sprintf("/%s", path)]; ok {
			currentResource = resource
			continue
		}

		//color.Green("%+v", s.resources)


	}
	return currentAction
}

func (s *TyphoonServer) requestHandler(ctx *gin.Context)  {
	logger := ginlogrus.GetCtxLogger(ctx)
	action := s.getAction(ctx)
	if action == nil {
		logger.Error("Not found")
		ctx.JSON(404, gin.H{
			"message": "Not Found",
			"status": false,
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"controller": action.Name,
	})

	controller := action.Controller

	controller(logger, ctx)

}

func (s *TyphoonServer) initActions(resource *interfaces.Resource)  {
	for name, action := range resource.Actions {
		for _, method := range action.Methods {
			var handlerPath string
			if resource.Path != "/" {
				handlerPath = fmt.Sprintf("%s/%s",resource.Path, name)
			} else {
				handlerPath = fmt.Sprintf("/%s", name)
			}

			//color.Yellow("init action %s ==> %s", name, handlerPath)

			s.Serve(method, handlerPath, s.requestHandler)

		}
	}
}

func (s *TyphoonServer) initResource(newResource *interfaces.Resource) error {
	if _, ok := s.resources[newResource.Path]; ok {
		return Errors.ResourceAlreadyExist
	} else {

		s.resources[newResource.Path] = newResource
		s.initActions(newResource)
	}
	return nil
}

func (s *TyphoonServer) resourcesServe(method string, path string, callback func(ctx *gin.Context))  {

	var handler func(ctx *gin.Context)

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

func (s *TyphoonServer) Serve(method string, path string, callback func(ctx *gin.Context))  {
	s.Init()
	if len(s.resources) == 0 {
		s.callbacks[path] = callback
		s.resourcesServe(method, path, callback)
	} else {
		s.resourcesServe(method, path, nil)
	}

}

func (s *TyphoonServer) CreateResource(path string, opts interfaces.BaseServerLabel) (error, *interfaces.Resource) {
	newResource := &interfaces.Resource{
		Path: path,
		Name: opts.Name,
		Description: opts.Description,
		Middlewares:     make([]*interfaces.Middleware, 0),
		Actions:         make(map[string]*interfaces.Action, 0),
	}
	err := s.initResource(newResource)
	return err, newResource
}

func (s *TyphoonServer) AddResource(resource *interfaces.Resource) error  {
	s.Init()
	err := s.initResource(resource)
	return err
}