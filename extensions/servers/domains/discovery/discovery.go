package discovery

import (
	v1 "github.com/vortex14/gotyphoon/extensions/servers/domains/discovery/resources/v1"
	"github.com/vortex14/gotyphoon/interfaces"
	serverInterface "github.com/vortex14/gotyphoon/interfaces/server"
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

) serverInterface.Interface {

	discoveryServer := (
		&server.TyphoonServer{
			Port: port,
			Level: interfaces.INFO,
			BaseServerLabel: &serverInterface.BaseServerLabel{
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
