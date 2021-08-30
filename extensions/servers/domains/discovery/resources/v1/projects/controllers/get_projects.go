package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/interfaces"
	"time"
)

const (

	ProjectsControllerName = "get_projects"
	ProjectsControllerDescription = "Get registered projects in discovery service"

)

type Component struct {
	MooredAt time.Time `json:"moored_at"`
	Port int `json:"port" binding:"required"`
	Host string `json:"host" binding:"required"`
	Cluster string `json:"cluster" binding:"required"`
	Component string `json:"component" binding:"required"`
}

type Project struct {
	Project string `json:"project" binding:"required"`
	*Component
}

var Projects = make(map[string]map[string]*Component)


// projectsHandler
// @Tags Auth
// @Accept  json
// @Produce  json
// @Summary Discovery projects controller
// @Description Typhoon Discovery projects controller
// @Success 200 {object} Projects
// @Router /api/v1/projects/get_projects [get]
func projectsHandler (ctx *gin.Context, logger interfaces.LoggerInterface ) {
	ctx.JSON(200, Projects)
}

var ProjectsController = &interfaces.Action{
	Name: ProjectsControllerName,
	Description: ProjectsControllerDescription,
	Controller: projectsHandler,
	Methods : []string{interfaces.GET},
}
