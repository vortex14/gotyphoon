package gin

import (
	Context "context"
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/ctx"
)

const (
	CTX = "GIN_CTX"
)

func NewRequestCtx(context Context.Context, ginCtx *gin.Context) Context.Context{
	return ctx.Update(context, CTX, ginCtx)
}

func GetRequestCtx(context Context.Context) (bool, *gin.Context){
	request, ok := ctx.Get(context, CTX).(*gin.Context)
	return ok, request
}
