package component

import (
	"github.com/vortex14/gotyphoon/interfaces"
	serverInterface "github.com/vortex14/gotyphoon/interfaces/server"
	"github.com/vortex14/gotyphoon/server"
)


func Constructor(
	name string,
	labels *serverInterface.BaseServerLabel,
	project interfaces.Project,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *interfaces.BaseLoggerOptions,
	swaggerOptions *interfaces.SwaggerOptions,

) serverInterface.Interface {

	componentServer := (
		&server.TyphoonServer{
			Port: project.LoadConfig().GetComponentPort(name),
			Level: project.GetLogLevel(),
			BaseServerLabel: labels,

			TracingOptions: tracingOptions,
			LoggerOptions: loggerOptions,
			SwaggerOptions: swaggerOptions,
		}).
		Init().
		InitLogger().
		InitTracer()

	return componentServer
}
