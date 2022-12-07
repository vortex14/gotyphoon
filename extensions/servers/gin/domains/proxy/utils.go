package proxy

import (
	b64 "encoding/base64"
	"net/url"

	Gin "github.com/gin-gonic/gin"
)

func GetUrlParamGin(ctx *Gin.Context) *url.URL {
	params := ctx.Request.URL.Query()
	sURL := params["url"]
	encodeFlag := params["encode"]
	isBase64 := false

	if len(encodeFlag) > 0 && encodeFlag[0] == "base64" {
		isBase64 = true
	}
	URL := sURL[0]
	if isBase64 {

		sDec, _ := b64.StdEncoding.DecodeString(URL)
		URL = string(sDec)
	}

	u, _ := url.Parse(URL)
	return u
}
