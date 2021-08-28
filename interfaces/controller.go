package interfaces

import (
	"github.com/gin-gonic/gin"
)

type Controller func(ctx *gin.Context, logger LoggerInterface)



