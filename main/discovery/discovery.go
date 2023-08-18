package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/discovery"
)

func main() {
	logrus.Info("starting discovery local server ...")
	_ = discovery.Constructor("localhost", 12735, "http",
		nil,
		nil,
		nil).Run()
}
