package log

import (
	"fmt"
	"io"
	"os"

	ginOpentracing "github.com/Bose/go-gin-opentracing"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"

	"github.com/fatih/color"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"

	"github.com/vortex14/gotyphoon/interfaces"
)

type Options struct {
	*BaseOptions
}

func (o *Options) InitFormatter()  {
	color.Yellow("Init Log Formatter ... line: %t, file: %t, short: %t", o.ShowLine, o.ShowFile, o.ShortFileName)
	formatter := runtime.Formatter{ChildFormatter: &logrus.TextFormatter{
		FullTimestamp: o.FullTimestamp,
	}}
	formatter.Line = o.ShowLine
	formatter.File = o.ShowFile
	formatter.BaseNameOnly = o.ShortFileName
	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)
	color.Yellow("Set Log Level: %s", o.Level)
	logrus.SetLevel(o.GetLevel(o.Level))
}

type TyphoonLogger struct {
	Name string
	Options
	closer io.Closer
	reporter jaeger.Reporter
	tracer opentracing.Tracer
	TracingOptions *interfaces.TracingOptions
}

func (l *TyphoonLogger) InitTracer()  {
	if l.TracingOptions == nil {
		return
	}
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}

	tracer, reporter, closer, err := ginOpentracing.InitTracing(
		fmt.Sprintf("Typhoon-%s:%s", l.Name, hostName), // service name for the traces
		l.TracingOptions.GetEndpoint(),                        // where to send the spans
		ginOpentracing.WithEnableInfoLog(false)) // WithEnableLogInfo(false) will not log info on every span sent... if set to true it will log and they won't be aggregated
	if err != nil {
		color.Red("%s", err.Error())
		panic("unable to init tracing")
	}

	l.tracer = tracer
	l.reporter = reporter
	l.closer = closer


	opentracing.SetGlobalTracer(tracer)
}

func (l *TyphoonLogger) GetTracerHeader() string {
	return fmt.Sprintf("api-request-%s-",l.Name)
}

func (l *TyphoonLogger) Init() *TyphoonLogger {
	l.InitFormatter()
	l.InitTracer()
	return l
}

func (l *TyphoonLogger) Stop()  {
	defer l.closer.Close()
	defer l.reporter.Close()
}