package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces/server"
	"net/http"
	"time"
)

const (
	RegisterControllerName = "register"
	RegisterControllerDescription = "Registering projects controller for Discovery server"
)

type Status struct {
	Status bool
}


// registerHandler
// @Tags Auth
// @Accept  json
// @Produce  json
// @Summary Registering projects controller
// @Description Typhoon registering projects controller
// @Success 200 {object} Status
// @Router /api/v1/projects/register [post]
func registerHandler (logger *logrus.Entry, ctx *gin.Context ) {
	var project *Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": Errors.BadRequest.Error(),
		})
		return
	}
	project.MooredAt = time.Now().UTC()
	logger.WithFields(logrus.Fields{
		"Host": project.Host,
		"Port": project.Port,
		"Cluster": project.Cluster,
		"Project": project.Project,
		"Component": project.Component.Component,
	}).Info("Registering of Typhoon component")

	component := &Component{
		MooredAt:  project.MooredAt,
		Port:      project.Port,
		Host:      project.Host,
		Cluster:   project.Cluster,
		Component: project.Component.Component,
	}

	if regProject := Projects[project.Project]; regProject == nil {
		Projects[project.Project] = map[string]*Component{
			project.Component.Component: component,
		}
	} else {
		Projects[project.Project][project.Component.Component] = component
	}

	ctx.JSON(200, &Status{Status: true})
}

var RegisterController = &server.Action{
	Name: RegisterControllerName,
	Description: RegisterControllerDescription,
	Controller: registerHandler,
	Methods : []string{server.POST},
}
