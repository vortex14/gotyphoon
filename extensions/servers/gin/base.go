package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

func SetServeHandler(
	method string,
	path string,
	group *gin.RouterGroup,
	handler func(ctx *gin.Context),
) {

	switch method {
	case interfaces.GET     : group.GET(path, handler)
	case interfaces.PUT     : group.PUT(path, handler)
	case interfaces.POST    : group.POST(path, handler)
	case interfaces.PATCH   : group.PATCH(path, handler)
	case interfaces.DELETE  : group.DELETE(path, handler)
	case interfaces.OPTIONS : group.OPTIONS(path, handler)
	}
}
