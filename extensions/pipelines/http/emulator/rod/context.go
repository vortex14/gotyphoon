package rod

import (
	Context "context"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/vortex14/gotyphoon/ctx"
)

const (
	BROWSER  = "browser"
	PAGE     = "page"
	LAUNCHER = "launcher"
	RESPONSE = "response"
)

func NewBrowserCtx(context Context.Context, browser *rod.Browser) Context.Context {
	return ctx.Update(context, BROWSER, browser)
}

func NewBodyResponse(context Context.Context, body *string) Context.Context {
	return ctx.Update(context, RESPONSE, body)
}

func GetPageResponse(context Context.Context) (bool, *string) {
	body, ok := ctx.Get(context, RESPONSE).(*string)
	return ok, body
}

func GetBrowserCtx(context Context.Context) (bool, *rod.Browser) {
	browser, ok := ctx.Get(context, BROWSER).(*rod.Browser)
	return ok, browser
}

func NewPageCtx(context Context.Context, page *rod.Page) Context.Context {
	return ctx.Update(context, PAGE, page)
}

func GetPageCtx(context Context.Context) (bool, *rod.Page) {
	page, ok := ctx.Get(context, PAGE).(*rod.Page)
	return ok, page
}

func NewLauncherCtx(context Context.Context, launcher *launcher.Launcher) Context.Context {
	return ctx.Update(context, LAUNCHER, launcher)
}

func GetLauncherCtx(context Context.Context) (bool, *launcher.Launcher) {
	launcher, ok := ctx.Get(context, LAUNCHER).(*launcher.Launcher)
	return ok, launcher
}
