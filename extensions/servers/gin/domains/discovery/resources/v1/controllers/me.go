package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"

	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
)

type MeParams struct {
	Id string `form:"id" binding:"required" json:"id"`
	//Name string `uri:"name" binding:"required"`
}

const (
	PATH                    = "me"
	NAMEMeController        = "me"
	DESCRIPTIONMeController = "Me discovery JWT controller extension for Typhoon server"

	JWTAUTHLOGINDefault = "typhoon"
	JWTAUTHEmailDefault = "typhoon@typhoon-s1.ru"
	JWTRoleDefault      = "admin"
)

type MeResponse struct {
	Login string   `json:"login"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// handler
// @Tags Auth
// @Accept  json
// @Produce  plain
// @Summary Discovery me controller
// @Description Typhoon Discovery me controller
// @Success 200 {object} MeResponse
// @Router /api/v1/me [get]
func meHandler(ctx *gin.Context, logger interfaces.LoggerInterface) {
	params := &MeParams{}

	if err := ctx.ShouldBindQuery(params); err != nil {

		//m := make(map[string][]string)
		//for _, v := range ctx. {
		//	m[v.Key] = []string{v.Value}
		//}
		//
		//println(fmt.Sprintf("%+v; %+v", m, ctx.Params, ctx.Quer))

		ctx.JSON(400, &forms.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(200, &MeResponse{
		Login: JWTAUTHLOGINDefault,
		Email: JWTAUTHEmailDefault,
		Roles: []string{JWTRoleDefault},
	})
}

var MeController = &GinExtension.Action{
	Action: &forms.Action{
		Params: &MeParams{},
		MetaInfo: &label.MetaInfo{
			Path:        PATH,
			Name:        NAMEMeController,
			Description: DESCRIPTIONMeController,
			Tags:        []string{"Auth"},
		},
		Methods: []string{interfaces.GET},
	},
	GinController: meHandler,
}
