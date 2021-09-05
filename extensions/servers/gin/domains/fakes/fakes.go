package fakes

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

const (
	NAME = "Fakes"
	DESCRIPTION = "Server for data fakes"

	FakeProductPath = "product"
	FakeTaskPath = "task"
	FakeProxyPath = "proxy"
	FakeImagePath = "image"

	WATERMARK = "image.typhoon.dev"
)

func Constructor(
	port int,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *log.Options,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {
	return (
		&gin.TyphoonGinServer{
			TyphoonServer: &forms.TyphoonServer{
				Port: port,
				Level: interfaces.INFO,
				MetaInfo: &label.MetaInfo{
					Name:        NAME,
					Description: DESCRIPTION,
				},
				TracingOptions: tracingOptions,
				LoggerOptions: loggerOptions,
				SwaggerOptions: swaggerOptions,
			},
		}).
		Init().
		InitLogger().
		AddResource(
			&forms.Resource{
				MetaInfo: &label.MetaInfo{
					Path: "/fake",
					Name: "Main faker",
					Description: "Main Resource",
				},
				Actions: map[string]interfaces.ActionInterface{
					ping.PATH:       ping.Controller,
					FakeProductPath: CrateProductAction(),
					FakeTaskPath:    CreateTaskAction(),
					FakeProxyPath:   CreateProxyAction(),
					FakeImagePath:   CreateImageAction(),
				},
			},
		)
	}