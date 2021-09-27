package auth

import (
	"github.com/gin-gonic/gin"
	ginExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type Users map[string]string

type BasicAuth struct {
	Users  Users
	server *gin.Engine
	LOG interfaces.LoggerInterface
}

func (a *BasicAuth) Allow(server interfaces.ServerInterface, resource interfaces.ResourceInterface) interface{} {
	authorizedGroup := a.server.Group(resource.GetPath(), gin.BasicAuth(gin.Accounts(a.Users)))

	server.SetRouterGroup(resource, authorizedGroup)

	return authorizedGroup
}

func (a *BasicAuth) SetLogger(logger interfaces.LoggerInterface)  {
	a.LOG = log.PatchLogI(logger, log.D{"auth": "basic-auth"})
}

func (a *BasicAuth) SetServerEngine(server interfaces.ServerInterface)  {

	ok, ginServer := ginExtension.GetTyphoonGinServer(server)
	if !ok { a.LOG.Error(Errors.BasicAuthContextFailed.Error()); return }

	ginEngine := ginExtension.GetGinEngine(ginServer)
	a.server = ginEngine
}