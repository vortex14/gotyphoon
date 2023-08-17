package discovery

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/discovery/resources/v1"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

const (
	NAME        = "Discovery"
	DESCRIPTION = "Project discovery service"
)

func init() {
	log.InitD()
}

func Constructor(
	host string,
	port int,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *log.Options,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {

	discoveryServer := (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			Host:  host,
			Port:  port,
			Level: interfaces.INFO,
			MetaInfo: &label.MetaInfo{
				Name:        NAME,
				Description: DESCRIPTION,
				Version:     "1.0.1",
			},

			TracingOptions: tracingOptions,
			LoggerOptions:  loggerOptions,
			SwaggerOptions: swaggerOptions,
		},
	}).
		Init().
		InitLogger().
		AddResource(v1.Constructor().Get())

	return discoveryServer
}
