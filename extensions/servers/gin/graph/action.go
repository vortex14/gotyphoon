package graph

// /* ignore for building amd64-linux
//
//import (
//	"context"
//
//	"github.com/gin-gonic/gin"
//
//	Errors "github.com/vortex14/gotyphoon/errors"
//	graphForms "github.com/vortex14/gotyphoon/extensions/forms/graph"
//	ginExt "github.com/vortex14/gotyphoon/extensions/servers/gin"
//	"github.com/vortex14/gotyphoon/interfaces"
//	"github.com/vortex14/gotyphoon/utils"
//)
//
//type ServerGraphController func(
//	ctx *gin.Context,
//	server interfaces.ServerGraphInterface,
//	logger interfaces.LoggerInterface,
//)
//
//
//
//type Action struct {
//	*graphForms.Action
//
//	GinController  ginExt.Controller
//	GinSController ServerGraphController
//}
//
//func (a *Action) OnRequest(method string, path string) {
//	println(":::>>",method, path)
//}
//
//func (a *Action) AddMethod(name string) {
//	switch name {
//	case interfaces.POST, interfaces.GET, interfaces.PUT, interfaces.PATCH, interfaces.DELETE: a.Methods = append(a.Methods, name)}
//}
//
//func (a *Action) Run(context context.Context, logger interfaces.LoggerInterface) {
//	if utils.IsNill(a.GinController, a.Pipeline) { logger.Error(Errors.ActionMethodsNotFound.Error()); return }
//	status, requestCtx := ginExt.GetRequestCtx(context)
//
//	if !status { logger.Error(Errors.ActionContextRequestFailed.Error()) }
//	if a.GinController != nil { a.GinController(requestCtx, logger) } else if a.Pipeline != nil {
//		a.Pipeline.Run(context)
//	}
//
//	if ok, server := ginExt.GetServerCtx(context); ok && a.GinSController != nil {
//		a.GinSController(requestCtx, server.(interfaces.ServerGraphInterface), logger )
//	}
//}
//
// */


