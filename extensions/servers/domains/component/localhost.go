package component

import "github.com/vortex14/gotyphoon/interfaces"

func ConstructorLocalhostOptions(
	componentName string,

	) (*interfaces.BaseLoggerOptions, *interfaces.TracingOptions, *interfaces.SwaggerOptions) {

		return &interfaces.BaseLoggerOptions{
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
