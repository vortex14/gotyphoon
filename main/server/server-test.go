package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/logger"
	"github.com/vortex14/gotyphoon/interfaces"
	TS "github.com/vortex14/gotyphoon/server"
)



func init()  {
	(&logger.TyphoonLogger{
		Name: "App",
		Options: logger.Options{
			BaseLoggerOptions: &interfaces.BaseLoggerOptions{
				Name:          "Test-App",
				Level:         "DEBUG",
				ShowLine:      true,
				ShowFile:      true,
				ShortFileName: true,
				FullTimestamp: true,
			},
		},
	}).Init()
}



func main()  {
	logrus.Debug("start new test server")

	server := (&TS.TyphoonServer{
		Port: 6666,
		IsDebug: true,

	}).
		Init().
		InitLogger()


	home := &forms.Resource{
		Path: "/",
		Name: "Home",
		Description: "Home Resource",
		Middlewares: []interfaces.MiddlewareInterface{
			&forms.Middleware{
				Required: true,
				Name:        "Resource middleware 1",
				Fn: func(context context.Context, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

				},
			},
		},
		Resources:  make(map[string]interfaces.ResourceInterface),
		Actions: map[string]interfaces.ActionInterface{
			"some-test": &interfaces.Action{
				Name:        "Action name",
				Description: "Action description",
				Path :       "some-test",
				Methods:     []string{interfaces.GET},
				Middlewares: []interfaces.MiddlewareInterface{
					&forms.Middleware{
						Required: true,
						Name:        "middleware-1",
						Fn: func(context context.Context, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

						},
					},
				},
				Controller:  func(ctx *gin.Context, logger interfaces.LoggerInterface) {
					logger.Debug("request in controller")


					ctx.JSON(200, gin.H{
						"test": 2,
					})
				},
			},
		},
	}


	server.AddResource(home)
	//u := utils.Utils{}
	//println(u.PrintPrettyJson(home))
	err := server.Run()

	if err != nil {
		logrus.Error(err.Error())
	}
}
