package fakes

import (
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

	FakeProductPath = "product"
	FakeTaskPath = "task"
	FakeProxyPath = "proxy"
)

func Constructor(
	port int,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *interfaces.BaseLoggerOptions,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {
	return (
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
				Path: "/fake",
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
							ctx.JSON(200, fake.CreateProduct())
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
					FakeProxyPath: &interfaces.Action{
						Name: NAME,
						Path: FakeProxyPath,
						Description: "Fake Typhoon proxy",
						Controller: func(ctx *gin.Context, logger interfaces.LoggerInterface) {
							ctx.JSON(200, fake.CreateFakeProxy())
						},
						Methods : []string{interfaces.GET},
					},
				},
			},
		)
	}
