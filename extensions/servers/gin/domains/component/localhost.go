package component

import (
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func ConstructorLocalhostOptions(
	componentName string,

	) (*log.BaseOptions, *interfaces.TracingOptions, *interfaces.SwaggerOptions) {

		return &log.BaseOptions{
			ShowLine: true,
			ShowFile: true,
			FullTimestamp: true,
			Name: componentName,
			Level: interfaces.DEBUG,
		},
		&interfaces.TracingOptions{
			UseUTC:        true,
			UseBanner:     false,
			EnableInfoLog: false,
			JaegerPort:    LocalhostJaegerPort,
			JaegerHost:    LocalhostJaegerHost,
		},
		&interfaces.SwaggerOptions{
			DocEndpoint: LocalhostSwaggerEndpointDefinition,
		}
}
