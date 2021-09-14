package main

import (
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/fakes"
	"github.com/vortex14/gotyphoon/log"
)

func init()  {
	log.InitD()
}


func main()  {
	logrus.Info("starting fakes server ...")
	_ = fakes.GraphConstructor(12666,
		nil,
		nil,
		nil).Run()
}
