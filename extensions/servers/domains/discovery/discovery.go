package discovery

import (
	v1 "github.com/vortex14/gotyphoon/extensions/servers/domains/discovery/resources/v1"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/server"
)

const (
	NAME = "Discovery"
	DESCRIPTION = "Project discovery service"
)

func Constructor(
	port int,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *interfaces.BaseLoggerOptions,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {

	discoveryServer := (
		&server.TyphoonServer{
			Port: port,
			Level: interfaces.INFO,
			BaseServerLabel: &interfaces.BaseServerLabel{
				Name: NAME,
				Description: DESCRIPTION,
			},
			TracingOptions: tracingOptions,
			LoggerOptions: loggerOptions,
			SwaggerOptions: swaggerOptions,
		}).
		Init().
		InitLogger().
		AddResource(v1.Constructor().Get())

	return discoveryServer
}
