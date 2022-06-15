package main

import (
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

func main() {

	//chrome := &ChromeHeadless{
	//	//Proxy: "http://154.9.51.35:8800",
	//}
	//
	//res := chrome.Request("https://www.samsclub.com/api/node/vivaldi/v1/category/6930116")
	////var proxy models.Proxy
	////_ = utils.JsonLoad(&proxy, res)
	//log.Println(res)

	mainId := "100001"

	url := fmt.Sprintf("https://www.samsclub.com/api/node/vivaldi/v1/category/%s", mainId)

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			fmt.Println(string(r.Body))
		},
		//BrowserEndpoint: "ws://localhost:3000",
	}).Start()

}