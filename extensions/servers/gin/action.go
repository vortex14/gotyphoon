package gin

import (
	"context"

	"github.com/vortex14/gotyphoon/elements/forms"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

type Action struct {
	*forms.Action

	GinController  Controller
	GinSController ServerController
	RequestModel   interface{}
}

func (a *Action) AddMethod(name string) {
	switch name {
	case interfaces.POST, interfaces.GET, interfaces.PUT, interfaces.PATCH, interfaces.DELETE:
		a.Methods = append(a.Methods, name)
	}
}

func (a *Action) Cancel(ctx context.Context, logger interfaces.LoggerInterface, err error) {
	if a.Cn != nil {
		a.Cn(err, ctx, logger)
	}
}

func (a *Action) Run(context context.Context, logger interfaces.LoggerInterface) {
	if utils.IsNill(a.GinController, a.Pipeline) {
		logger.Error(Errors.ActionMethodsNotFound.Error())
		return
	}

	status, requestCtx := GetRequestCtx(context)
	if !status {
		logger.Error(Errors.ActionContextRequestFailed.Error())
	}
	if a.GinController != nil {
		a.GinController(requestCtx, logger)
	} else if a.Pipeline != nil {
		a.SafeRun(func() error {
			return a.Pipeline.Run(context)
		}, func(err error) {
			a.Cancel(context, logger, err)
		})

	}
	if ok, server := GetServerCtx(context); ok && a.GinSController != nil {
		a.GinSController(requestCtx, server, logger)
	}
}
