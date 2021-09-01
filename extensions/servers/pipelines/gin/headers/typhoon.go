package headers

import (
	"context"
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/pipelines/gin"
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
	request *Gin.Context,
) error {



	return nil

}

func (h *HttpRequestSetHeaderMiddleware) Pass(
	context context.Context,
	loggerInterface interfaces.LoggerInterface,
	reject func(err error),
	next func(ctx context.Context),
	) {

	request, _ := context.Value(gin.REQUEST).(*Gin.Context)
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
