package forms

import (
	Context "context"
	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/models/bar"
	"github.com/vortex14/gotyphoon/elements/models/label"
)

const (
	GOTOStage = "GOTO_STAGE"
	PlABEL    = "PIPELINE_LABEL"
	PrBAR     = "PROGRESS_BAR"
)

func PatchCtxPipelineGOTO(context Context.Context, stage int) Context.Context {
	return ctx.Update(context, GOTOStage, stage)
}

func GetGOTOCtx(context Context.Context) (bool, int) {
	stageIndex, ok := ctx.Get(context, GOTOStage).(int)
	return ok, stageIndex
}

func setLabel(context Context.Context, label *label.MetaInfo) Context.Context {
	return ctx.Update(context, PlABEL, label)
}

func GetPipelineLabel(context Context.Context) (bool, *label.MetaInfo) {
	_label, ok := ctx.Get(context, PlABEL).(*label.MetaInfo)
	return ok, _label
}

func setBar(context Context.Context, label *bar.Bar) Context.Context {
	return ctx.Update(context, PrBAR, label)
}
func GetBar(context Context.Context) (bool, *bar.Bar) {
	_bar, ok := ctx.Get(context, PrBAR).(*bar.Bar)
	return ok, _bar
}
