package gin

import (
	Context "context"
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	CTX = "GIN_CTX"
	TYPHOONServer = "TYPHOON_SERVER"
)


func GetTyphoonGinServer(server interfaces.ServerInterface) (bool, *TyphoonGinServer){
	ginServer, ok := server.(*TyphoonGinServer)
	return ok, ginServer
}

func GetGinEngine(server *TyphoonGinServer) *gin.Engine {
	return server.GetServerEngine().(*gin.Engine)
}

func GetGinEngineI(server interface{}) *gin.Engine {
	return server.(*gin.Engine)
}

func GetGinGroup(group interface{}) *gin.RouterGroup {
	return group.(*gin.RouterGroup)
}

func NewRequestCtx(context Context.Context, ginCtx *gin.Context) Context.Context{
	return ctx.Update(context, CTX, ginCtx)
}

func GetRequestCtx(context Context.Context) (bool, *gin.Context){
	request, ok := ctx.Get(context, CTX).(*gin.Context)
	return ok, request
}

func GetServerCtx(context Context.Context) (bool, interfaces.ServerInterface){
	request, ok := ctx.Get(context, TYPHOONServer).(interfaces.ServerInterface)
	return ok, request
}

func NewServerCtx(context Context.Context, server interfaces.ServerInterface) Context.Context{
	return ctx.Update(context, TYPHOONServer, server)
}