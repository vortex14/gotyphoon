package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/log"
)

func init()  {
	log.InitD()
}


func main()  {
	logrus.Debug("checking middlewares. pre tests")

	//
	//fakeTask, _ := fake.CreateFakeTask(interfaces.FakeTaskOptions{
	//	UserAgent:   false,
	//	Cookies:     false,
	//	Auth:        false,
	//	Proxy:       false,
	//	AllowedHttp: nil,
	//})

	//fakeTask.Fetcher.Auth = map[string]string{
	//	"password": "sUwF}r#LXcly8%U5",
	//	"login": "Esom",
	//}

	//fakeTask.Fetcher.IsProxyRequired = true
	//fakeTask.URL = "https://httpstat.us/200"
	//
	//pipeline := httpPipeline.Constructor(fakeTask, nil)
	//
	//err, _ := pipeline.Run(context.TODO())
	//if err != nil {
	//	logrus.Error(err.Error())
	//}
}
