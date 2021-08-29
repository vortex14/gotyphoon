package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"net/http"
)

const (
	UnmoorControllerName = "unmoor"
	unmoorControllerDescription = "Unmoor projects from discovery"
)


type unmoorRequest struct {
	Project string `json:"project"`
}

// unmoorHandler
// @Tags Auth
// @Accept  json
// @Produce  json
// @Summary unmoor the project
// @Description Typhoon unmoor projects controller
// @Success 200 {object} Status
// @Router /api/v1/projects/unmoor [post]
func unmoorHandler (ctx *gin.Context, logger interfaces.LoggerInterface ) {
	var project *unmoorRequest
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": Errors.BadRequest.Error(),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"Project": project.Project,
	}).Info("to unmoor project ...")


	if regProject := Projects[project.Project]; regProject != nil {
		delete(Projects, project.Project)
	}

	ctx.JSON(200, &Status{Status: true})
}

var UnmoorController = &interfaces.Action{
	Name: UnmoorControllerName,
	Description: unmoorControllerDescription,
	Controller: unmoorHandler,
	Methods : []string{interfaces.POST},
}
