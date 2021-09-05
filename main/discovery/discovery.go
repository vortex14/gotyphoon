package main

import (
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/discovery"
)

func main()  {
	println("start discovery")
	_ = discovery.Constructor(12735,
		nil,
		nil,
		nil).Run()
}