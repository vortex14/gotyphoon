package forms

import (
	Context "context"
	"github.com/vortex14/gotyphoon/ctx"
)

const (
	GOTOStage = "GOTO_STAGE"
)

func PatchCtxPipelineGOTO(context Context.Context, stage int) Context.Context {
	return ctx.Update(context, GOTOStage, stage)
}

func GetGOTOCtx(context Context.Context) (bool, int) {
	stageIndex, ok := ctx.Get(context, GOTOStage).(int)
	return ok, stageIndex
}
