package server

import "github.com/gin-gonic/gin"

type Middleware struct {
	Name string
	Description string
	ctx *gin.Context
	Callback func(ctx *gin.Context)
	PyCallback func(ctx *gin.Context)
}

type MiddlewareInterface interface {
	MetaDataInterface
}

