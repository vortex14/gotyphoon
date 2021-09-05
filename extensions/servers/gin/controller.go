package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Controller func(ctx *gin.Context, logger interfaces.LoggerInterface)
