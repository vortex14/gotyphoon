package component

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func Constructor(
	name string,
	labels *label.MetaInfo,
	project interfaces.Project,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *log.Options,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {

	componentServer := (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			MetaInfo: labels,
			Port:     project.LoadConfig().GetComponentPort(name),
			Level:    project.GetLogLevel(),

			TracingOptions: tracingOptions,
			LoggerOptions:  loggerOptions,
		},
	}).
		Init().
		InitLogger().
		InitTracer()

	return componentServer
}
