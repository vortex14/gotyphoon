package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/discovery"
)

func main()  {
	logrus.Info("starting discovery local server ...")
	_ = discovery.Constructor(12735,
		nil,
		nil,
		nil).Run()
}