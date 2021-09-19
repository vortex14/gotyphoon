package graph

// /* ignore for building amd64-linux

import (
	"fmt"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"reflect"
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	ginopentracing "github.com/Bose/go-gin-opentracing"
	Gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/ctx"
	Errors "github.com/vortex14/gotyphoon/errors"
	graphFormExt "github.com/vortex14/gotyphoon/extensions/forms/graph"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"
)

type ServerExtension struct {

}

func (e *ServerExtension) onInitAction(resource interfaces.ResourceInterface, action interfaces.ActionInterface)  {

}

type TyphoonGraphGinServer struct {
	*graphFormExt.TyphoonServer

	server *Gin.Engine
}

func (s *TyphoonGraphGinServer) onInitAction(resource interfaces.ResourceInterface, action interfaces.ActionInterface) {
	s.LOG.Info("onInitAction path: ", action.GetHandlerPath(), reflect.TypeOf(action), reflect.TypeOf(resource))
	r, ro := resource.(interfaces.ResourceGraphInterface)
	a, ao := action.(interfaces.ActionGraphInterface)
	s.LOG.Error(r, a, ro ,ao)

	if graphAction, ok := action.(interfaces.ActionGraphInterface); ok {

		if graphResource, okR := resource.(interfaces.ResourceGraphInterface); okR {
			graphResource.AddGraphActionNode(graphAction)
		} else {
			s.LOG.Error(Errors.GraphActionContextInvalid.Error())
		}

	} else {
		s.LOG.Error(Errors.GraphActionContextInvalid.Error())
	}

}

func (s *TyphoonGraphGinServer) InitTracer() interfaces.ServerInterface {
	if utils.NotNill(s.Logger, s.TracingOptions) {
		p := ginopentracing.OpenTracer([]byte(s.Logger.GetTracerHeader()))

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

	} else { s.LOG.Error(Errors.TracerContextNotFound.Error()) }

	return s
}
//*/

// requestHandler handle all HTTP request in here

// /* ignore for building amd64-linux
func (s *TyphoonGraphGinServer) onRequestHandler(ginCtx *Gin.Context)  {

	requestContext := gin.NewRequestCtx(ctx.New(), ginCtx)
	requestLogger := ginlogrus.GetCtxLogger(ginCtx)

	reservedRequestPath := ginCtx.Request.URL.Path

	requestContext = s.InitRequestPath(requestContext, reservedRequestPath)

	action := s.GetAction(reservedRequestPath, requestLogger, ginCtx)
	action.OnRequest(ginCtx.Request.Method, reservedRequestPath)
	if action == nil { s.LOG.Error(Errors.ActionPathNotFound.Error())
		ginCtx.JSON(404, Gin.H{ "message": "Not Found", "status": false}); return
	}

	requestLogger = log.Patch(requestLogger, log.D{"controller": action.GetName()})
	requestContext = gin.NewServerCtx(requestContext, s)

	requestContext = log.NewCtx(requestContext, requestLogger)

	requestLogger.Debug(fmt.Sprintf("found action %s", action.GetName()))
	errStack, statusMiddlewareStack, _ := s.RunMiddlewareStack(requestContext, action)
	requestLogger.Debug(fmt.Sprintf("status middleware stack: %t", statusMiddlewareStack))

	if statusMiddlewareStack { action.Run(requestContext, requestLogger) } else {
		requestLogger.Debug(fmt.Sprintf("error middleware stack: %s", errStack.Error()))
	}

}

func (s *TyphoonGraphGinServer) onServeHandler(method string, path string)  {

	s.LOG.Debug(fmt.Sprintf("gin serve %s %s ",method, path))
	gin.SetServeHandler(method, path, s.server, s.onRequestHandler)
}

func (s *TyphoonGraphGinServer) OnStartGin(port int) error {
	s.LOG.Info(fmt.Sprintf("running server: %s : %d", s.GetName(), port))
	return s.server.Run(fmt.Sprintf(":%d", port))
}

func (s *TyphoonGraphGinServer) onInitResource(newResource interfaces.ResourceInterface)  {
	if graphResource, ok := newResource.(interfaces.ResourceGraphInterface); ok {
		s.LOG.Info("onInitResource, hasGraph: ", graphResource.HasParentGraph())
	}
}

// */

func (s *TyphoonGraphGinServer) onBuildSubResources(subResource interfaces.ResourceInterface)  {
	s.LOG.Warning("OnBuildSubResources")

	//subGraph := newResource.CreateSubGraph(&interfaces.GraphOptions{
	//	Name:      subResource.GetName(),
	//	Label:     subResource.GetName(),
	//	IsCluster: true,
	//})
	//subResource.SetGraph(subGraph)
	////
	//subResource.SetGraphNodes(newResource.GetGraphNodes())
	//subResource.SetGraphEdges(newResource.GetGraphEdges())
}

func (s *TyphoonGraphGinServer) onBuildSubAction(resource interfaces.ResourceInterface, action interfaces.ActionInterface)  {
	s.LOG.Info("onBuildSubAction")
}

// /* ignore for building amd64-linux

func (s *TyphoonGraphGinServer) onAddResource(resource interfaces.ResourceInterface)  {
	s.LOG.Info("onAddResource", resource)
	if graphResource, ok := resource.(interfaces.ResourceGraphInterface); ok {
		s.AddNewGraphResource(graphResource)
	} else {
		s.LOG.Error(Errors.GraphResourceContextInvalid.Error())
	}


}

func (s *TyphoonGraphGinServer) Init() interfaces.ServerInterface {

	s.Construct(func () {
		s.InitLogger()
		s.LOG.Debug("init Graph Typhoon Gin Server")
		s.InitResourcesMap()

		s.server = Gin.New()
		s.server.Use(Gin.Recovery())

		s.InitGraph()
		s.OnStart = s.OnStartGin
		s.OnServeHandler = s.onServeHandler
		s.OnBuildSubAction = s.onBuildSubAction
		s.OnBuildSubResources = s.onBuildSubResources
		s.OnInitResource = s.onInitResource
		s.OnAddResource = s.onAddResource
		s.OnInitAction = s.onInitAction

	})
	return s
}


// */