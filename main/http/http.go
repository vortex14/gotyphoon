package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/data/fake"
	"github.com/vortex14/gotyphoon/extensions/logger"
	httpPipeline "github.com/vortex14/gotyphoon/extensions/pipelines/http/strategies/default"
	"github.com/vortex14/gotyphoon/interfaces"
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
	logrus.Debug("checking middlewares. pre tests")


	fakeTask, _ := fake.CreateFakeTask(interfaces.FakeTaskOptions{
		UserAgent:   false,
		Cookies:     false,
		Auth:        false,
		Proxy:       false,
		AllowedHttp: nil,
	})

	//fakeTask.Fetcher.Auth = map[string]string{
	//	"password": "sUwF}r#LXcly8%U5",
	//	"login": "Esom",
	//}

	fakeTask.Fetcher.IsProxyRequired = true
	fakeTask.URL = "https://httpstat.us/200"

	pipeline := httpPipeline.Constructor(fakeTask, nil)

	err, _ := pipeline.Run(context.TODO())
	if err != nil {
		logrus.Error(err.Error())
	}
}
