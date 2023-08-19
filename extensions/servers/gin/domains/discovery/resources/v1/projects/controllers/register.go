package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"net/http"
	"time"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"

	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
)

const (
	RegisterControllerName        = "register"
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
func registerHandler(ctx *gin.Context, logger interfaces.LoggerInterface) {
	var project *Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": Errors.BadRequest.Error(),
		})
		return
	}
	project.MooredAt = time.Now().UTC()
	logger.WithFields(logrus.Fields{
		"Host":      project.Host,
		"Port":      project.Port,
		"Cluster":   project.Cluster,
		"Project":   project.Project,
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

var RegisterController = &GinExtension.Action{
	Action: &forms.Action{
		MetaInfo: &label.MetaInfo{
			Name:        RegisterControllerName,
			Description: RegisterControllerDescription,
			Tags:        []string{"Project"},
		},
		Methods: []string{interfaces.POST},
	},
	GinController: registerHandler,
}
