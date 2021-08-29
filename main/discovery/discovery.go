package main

import (
	"github.com/vortex14/gotyphoon/extensions/servers/domains/discovery"
)

func main()  {
	println("start discovery")
	_ = discovery.Constructor(12735,
		nil,
		nil,
		nil).Init().InitLogger().Run()
}