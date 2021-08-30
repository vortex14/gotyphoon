package fakes

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/data/fake"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/servers/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/server"
)

const (
	NAME = "Fakes"
	DESCRIPTION = "Server for data fakes"

	FakeProductPath = "product-fake"
	FakeTaskPath = "task-fake"
)

func Constructor(
	port int,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *interfaces.BaseLoggerOptions,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {

	FakerServer := (
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
		AddResource(
			&forms.Resource{
				Path: "/",
				Name: "Main faker",
				Description: "Main Resource",
				Resources:    make(map[string]interfaces.ResourceInterface),
				Middlewares: make([]interfaces.MiddlewareInterface, 0),
				Actions: map[string]interfaces.ActionInterface{
					ping.PATH: ping.Controller,
					FakeProductPath: &interfaces.Action{
						Name: NAME,
						Path: FakeProductPath,
						Description: "Fake Product",
						Controller: func(ctx *gin.Context, logger interfaces.LoggerInterface) {

							var f fake.Product
							err := gofakeit.Struct(&f)
							if err != nil {
								color.Red("%s", err.Error())
								ctx.JSON(400, gin.H{"status": false})
								return
							}

							ctx.JSON(200, f)
						},
						Methods : []string{interfaces.GET},
					},
					FakeTaskPath: &interfaces.Action{
						Name: NAME,
						Path: FakeTaskPath,
						Description: "Fake Typhoon task",
						Controller: func(ctx *gin.Context, logger interfaces.LoggerInterface) {
							ctx.JSON(200, fake.CreateDefaultTask())
						},
						Methods : []string{interfaces.GET},
					},
				},
			},
		)

	return FakerServer
}
