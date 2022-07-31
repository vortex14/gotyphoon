package main

import (
	_u "github.com/ahl5esoft/golang-underscore"
	"github.com/elazarl/goproxy"
	"github.com/fatih/color"
	"log"
	"net/http"
)

func main() {

	proxyList := []string{"1", "2", "3"}
	index := _u.Chain(proxyList).FindIndex(func(r string, _ int) bool {
		return "3" == r
	})
	color.Red("%d", index)
	return

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":11313", proxy))
}
