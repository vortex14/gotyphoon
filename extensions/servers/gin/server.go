package gin

import (
	"fmt"
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/Bose/go-gin-opentracing"
	Gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/forms"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"
)

type TyphoonGinServer struct {
	*forms.TyphoonServer

	server *Gin.Engine
}

func (s *TyphoonGinServer) InitTracer() interfaces.ServerInterface {
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

// requestHandler handle all HTTP request in here
func (s *TyphoonGinServer) onRequestHandler(ginCtx *Gin.Context)  {
	requestContext := NewRequestCtx(ctx.New(), ginCtx)
	requestLogger := ginlogrus.GetCtxLogger(ginCtx)

	reservedRequestPath := ginCtx.Request.URL.Path
	requestContext = s.InitRequestPath(requestContext, reservedRequestPath)

	action := s.GetAction(reservedRequestPath, requestLogger, ginCtx)

	if action == nil { s.LOG.Error(Errors.ActionPathNotFound.Error())
		ginCtx.JSON(404, Gin.H{
			"message": "Not Found",
			"status": false,
		}); return
	}

	requestLogger = log.Patch(requestLogger, log.D{"controller": action.GetName()})

	requestContext = log.NewCtx(requestContext, requestLogger)

	requestLogger.Debug(fmt.Sprintf("found action %s", action.GetName()))
	errStack, statusMiddlewareStack, _ := s.RunMiddlewareStack(requestContext, action)
	requestLogger.Debug(fmt.Sprintf("status middleware stack: %t", statusMiddlewareStack))
	if statusMiddlewareStack {
		action.Run(requestContext, requestLogger)
	} else {
		requestLogger.Debug(fmt.Sprintf("error middleware stack: %t", errStack.Error()))
	}

}

func (s *TyphoonGinServer) onServeHandler(method string, path string)  {

	s.LOG.Debug(fmt.Sprintf("gin serve %s %s ",method, path))
	switch method {
	case interfaces.GET:
		s.server.GET(path, s.onRequestHandler)
	case interfaces.POST:
		s.server.POST(path, s.onRequestHandler)
	case interfaces.PUT:
		s.server.PUT(path, s.onRequestHandler)
	case interfaces.PATCH:
		s.server.PATCH(path, s.onRequestHandler)
	case interfaces.DELETE:
		s.server.DELETE(path, s.onRequestHandler)
	}
}

func (s *TyphoonGinServer) onStart(port int) error {
	s.LOG.Info(fmt.Sprintf("running server: %s : %d", s.GetName(), port))
	return s.server.Run(fmt.Sprintf(":%d", port))
}

func (s *TyphoonGinServer) Init() interfaces.ServerInterface {

	s.Instance.Do(func () {
		s.InitLogger()
		s.LOG.Debug("init Typhoon Gin Server")
		s.InitResourcesMap()

		s.server = Gin.New()
		s.server.Use(Gin.Recovery())

		s.OnStart = s.onStart
		s.OnServeHandler = s.onServeHandler
	})
	return s
}

func (s *TyphoonGinServer) Stop() error  {
	return nil
}

func (s *TyphoonGinServer) Restart() error {
	return nil
}