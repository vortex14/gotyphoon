package gin

import (
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

const (
	LocalhostSwaggerEndpointDefinition = "http://localhost:8000/swagger/doc.json"
	LocalhostJaegerPort = 5775
	LocalhostJaegerHost = "localhost"
)

func ConstructorLocalhostOptions(
	componentName string,

	) (*log.Options, *interfaces.TracingOptions, *interfaces.SwaggerOptions) {

		return &log.Options{
			BaseOptions:&log.BaseOptions{
				ShowLine: true,
				ShowFile: true,
				FullTimestamp: true,
				Name: componentName,
				Level: interfaces.DEBUG,
			},

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
