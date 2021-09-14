package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

func SetServeHandler(
	method string,
	path string,
	server *gin.Engine,
	handler func(ctx *gin.Context),
) {

	switch method {
	case interfaces.GET    : server.GET(path, handler)
	case interfaces.PUT    : server.PUT(path, handler)
	case interfaces.POST   : server.POST(path, handler)
	case interfaces.PATCH  : server.PATCH(path, handler)
	case interfaces.DELETE : server.DELETE(path, handler)
	}
}
