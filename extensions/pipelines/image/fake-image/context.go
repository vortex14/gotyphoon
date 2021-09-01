package fake_image

import (
	Context "context"
	"github.com/fogleman/gg"
	"github.com/vortex14/gotyphoon/ctx"
)

func GetImgCtx(context Context.Context) (bool, *gg.Context){
	imgCtx, ok := ctx.Get(context, ImgCtx).(*gg.Context)
	return ok, imgCtx
}


func NewImgCtx(context Context.Context, data *gg.Context) Context.Context{
	return ctx.Update(context, ImgCtx, data)
}