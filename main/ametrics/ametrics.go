package main

import (
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/auth"
	"github.com/vortex14/gotyphoon/interfaces/server"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/graph"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/component"
	GinMiddlewares "github.com/vortex14/gotyphoon/extensions/servers/gin/middlewares"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/interfaces"
)

func main() {

	loggerOpt, tracingOpt, swaggerOpt := component.ConstructorLocalhostOptions("Ametrics")

	var mongoService = &mongo.Service{
		Settings: interfaces.ServiceMongo{
			Name: "Test",
			Details: struct {
				AuthSource string `yaml:"authSource,omitempty"`
				Username   string `yaml:"username,omitempty"`
				Password   string `yaml:"password,omitempty"`
				Host       string `yaml:"host"`
				Port       int    `yaml:"port"`
			}{
				Host: "localhost",
				Port: 27017,
			},
		},
	}

	aMetricServer := (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			IsDebug: true,
			Port:    14000,
			Level:   interfaces.DEBUG,
			MetaInfo: &label.MetaInfo{
				Name:        "Agent metric",
				Description: "Agent metric server",
			},

			TracingOptions: tracingOpt,
			LoggerOptions:  loggerOpt,
			SwaggerOptions: swaggerOpt,
		},
	}).
		Init().
		AddResource(
			&forms.Resource{
				MetaInfo: &label.MetaInfo{
					Path:        "/",
					Name:        "AResource",
					Description: "AResource process tracking",
				},
				Auth: []interfaces.ResourceAuthInterface{
					&auth.BasicAuth{
						Users: map[string]string{
							"vortex": "-=-=",
						},
					},
				},
				Actions: map[string]interfaces.ActionInterface{
					ping.PATH:  ping.Controller,
					graph.PATH: graph.Controller,

					"tracker": &GinExtension.Action{
						Action: &forms.Action{
							Middlewares: []interfaces.MiddlewareInterface{
								GinMiddlewares.ConstructorCorsMiddleware(server.GetAllAllowedCors()),
							},
							MetaInfo: &label.MetaInfo{
								Name:        "tracker",
								Path:        "tracker",
								Description: "listen new data track",
							},
							Methods: []string{interfaces.POST, interfaces.OPTIONS, interfaces.GET},
						},
						GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
							ctx.JSON(200, Gin.H{"ping": mongoService.Ping()})
						},
					},
				},
			},
		)

	err := aMetricServer.Run()
	if err != nil {
		return
	}
}
