package required

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/extensions/servers/middlewares/gin-gonic/headers"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	NAME = "Checking the header of a client sending a request to the server"
	DESCRIPTION = "Default header middleware wrapper for the Typhoon server that provided Gin context interface"
)

type HttpRequestSetHeaderMiddleware struct {
	*interfaces.BaseLabel
}

func (h *HttpRequestSetHeaderMiddleware) Run(
	request *gin.Context,
) error {



	return nil

}

func (h *HttpRequestSetHeaderMiddleware) Pass(
	context context.Context,
	loggerInterface interfaces.LoggerInterface,
	reject func(err error),

	) {

	request, _ := context.Value(headers.REQUEST).(*gin.Context)
	err := h.Run(request)
	if err != nil {
		reject(err)
	}
}

func ConstructorRequestHeaderMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpRequestSetHeaderMiddleware{
		BaseLabel: &interfaces.BaseLabel{
			Required:    required,
			Name:        NAME,
			Description: DESCRIPTION,
		},
	}
}
