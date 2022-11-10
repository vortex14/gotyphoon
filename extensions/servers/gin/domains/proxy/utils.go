package proxy

import (
	"net/url"

	Gin "github.com/gin-gonic/gin"
)

func GetUrlParamGin(ctx *Gin.Context) *url.URL {
	params := ctx.Request.URL.Query()
	sURL := params["url"]
	u, _ := url.Parse(sURL[0])
	return u
}
