package main

import (
	"github.com/vortex14/gotyphoon/extensions/servers/gin/domains/proxy"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func main() {
	_ = proxy.Constructor("localhost").Run()

}
