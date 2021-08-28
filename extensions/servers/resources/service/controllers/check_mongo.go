package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	CheckMongoPath = "check_mongo"
	CheckMongoDescription = "Health check controller of Mongo"

	LocalhostMongoHost = "localhost"
	localhostMongoPort = 27017

)

var mongoService = &mongo.Service{
	Settings: interfaces.ServiceMongo{
		Name: "Test",
		Details: struct {
    		AuthSource string `yaml:"authSource,omitempty"`
			Username   string `yaml:"username,omitempty"`
			Password   string `yaml:"password,omitempty"`
			Host       string `yaml:"host"`
			Port       int    `yaml:"port"`
		}{
			Host: LocalhostMongoHost,
			Port: localhostMongoPort,
		},
	},
}

// MongoHandler
// @Tags Services
// @Produce  json
// @Summary controller of Mongo
// @Description Health check controller of Mongo
// @Success 200 {object} ServiceResponse
// @Router /api/v1/services/check_mongo [get]
func MongoHandler (logger *logrus.Entry, ctx *gin.Context ) {
	ctx.JSON(200, &ServiceResponse{
		Status: mongoService.Ping(),
	})
}

var MongoController = &interfaces.Action{
	Name: CheckMongoPath,
	Description: CheckMongoDescription,
	Controller: MongoHandler,
	Methods : []string{interfaces.GET},
}

