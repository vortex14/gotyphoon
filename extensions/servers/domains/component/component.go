package component

import (
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/server"
)


func Constructor(
	name string,
	labels *interfaces.BaseServerLabel,
	project interfaces.Project,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *interfaces.BaseLoggerOptions,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {

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
