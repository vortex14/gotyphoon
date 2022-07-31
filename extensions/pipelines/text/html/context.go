package html

import (
	Context "context"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/ctx"
)

func GetHtmlDoc(context Context.Context) (bool, *goquery.Document) {
	htmlCtx, ok := ctx.Get(context, CtxHtml).(*goquery.Document)
	return ok, htmlCtx
}

func NewHtmlCtx(context Context.Context, data *goquery.Document) Context.Context {
	return ctx.Update(context, CtxHtml, data)
}
