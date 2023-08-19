package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"

	"github.com/vortex14/gotyphoon/interfaces"

	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
)

const (
	NAME        = "login"
	DESCRIPTION = "Discovery login controller extension for Typhoon server"

	JWTAUTHDefault = "eyJsb2dpbiI6InR5cGhvb24iLCJlbWFpbCI6InR5cGhvb25AdHlwaG9vbi1zMS5ydSIsInJvbGVzIjpbXSwiYWxnIjoiSFMy+" +
		"NTYifQ.e30.7m63q7oIzRooWceOw5DX-S8av4NHx_AbQx8oibISgZU"
)

type UserPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func handler(ctx *gin.Context, logger interfaces.LoggerInterface) {

	service, ok := ctx.Get(GinExtension.TYPHOONActionService)

	if ok {
		_service := service.(*Service)
		logger.Warning(fmt.Sprintf("OK ! %s", _service.Test()))
	}

	ctx.String(200, JWTAUTHDefault)
}

type Service struct {
	Repository interface{}
}

func (s *Service) Login() {

}

func (s *Service) Test() string {
	return "123"
}

var LoginController = &GinExtension.Action{
	Action: &forms.Action{
		BodyRequestModel: forms.BaseModelRequest{RequestModel: &UserPayload{}, Required: true},
		Service:          &Service{},
		MetaInfo: &label.MetaInfo{
			Tags:        []string{"Auth"},
			Path:        NAME,
			Name:        NAME,
			Description: DESCRIPTION,
		},
		ResponseModels: map[int]interface{}{
			200: &TokenResponse{},
		},
		Methods: []string{interfaces.POST},
	},
	GinController: handler,
}
