package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/fakes"
)

func main()  {
	logrus.Info("starting fakes server ...")
	_ = fakes.Constructor(12666,
		nil,
		nil,
		nil).Run()
}
