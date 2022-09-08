package json

import (
	Context "context"
	"github.com/itchyny/gojq"
	"github.com/vortex14/gotyphoon/ctx"
)

const (
	CtxJQ = "JQ_CTX"
)

func GetJQ(context Context.Context) (bool, gojq.Iter) {
	JQCtx, ok := ctx.Get(context, CtxJQ).(gojq.Iter)
	return ok, JQCtx
}

func NewJQCtx(context Context.Context, jq gojq.Iter) Context.Context {
	return ctx.Update(context, CtxJQ, jq)
}
