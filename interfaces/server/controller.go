package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Controller func(logger *logrus.Entry, ctx *gin.Context)



