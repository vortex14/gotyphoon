package proxy

import (
	Errors "errors"
	"fmt"
	u_ "github.com/ahl5esoft/golang-underscore"
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	netHttp "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/graph"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
)

func Constructor(opts *Settings) interfaces.ServerInterface {

	proxyCollection := Collection{
		Settings: opts,
	}

	proxyCollection.Init()

	return (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			Port:  opts.Port,
			Level: interfaces.INFO,
			MetaInfo: &label.MetaInfo{
				Name:        NAME,
				Description: DESCRIPTION,
			},
		},
	}).Init().InitLogger().AddResource(
		&forms.Resource{
			MetaInfo: &label.MetaInfo{
				Path:        PATH,
				Name:        ResourceName,
				Description: DESCRIPTION,
			},
			Actions: map[string]interfaces.ActionInterface{
				ping.PATH:  ping.Controller,
				graph.PATH: graph.Controller,
				"proxy": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "proxy",
							Path:        "proxy",
							Description: "Get a new proxy",
						},
						Methods: []string{interfaces.GET},
						Middlewares: []interfaces.MiddlewareInterface{
							netHttp.UrlRequiredMiddleware,
						},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get a new proxy.")
						proxy := fake.CreateFakeProxy()

						u := GetUrlParamGin(ctx)

						err, proxyValue := proxyCollection.GetProxy(u)
						if err != nil {
							logger.Debug(err.Error())
							proxy.Success = false
							proxy.Proxy = err.Error()
						}

						proxy.Proxy = proxyValue
						ctx.JSON(200, proxy)
					},
				},
				"blocked": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "proxy",
							Path:        "blocked",
							Description: "blocked proxies",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get stats of blocked proxies")

						ctx.JSON(200, proxyCollection.banned)
					},
				},
				"is_blocked": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "is blocked?",
							Path:        "is_blocked",
							Description: "is blocked the proxy ?",
						},
						Methods: []string{interfaces.GET},
						Middlewares: []interfaces.MiddlewareInterface{
							netHttp.UrlRequiredMiddleware,
						},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get stats of blocked proxies")

						params := ctx.Request.URL.Query()
						u := GetUrlParamGin(ctx)
						proxy := params["proxy"]
						if len(proxy) == 0 {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: Errors.New("query params hasn't url or proxy").Error()})
							return
						}

						index := u_.Chain(proxyCollection.banned[u.Hostname()]).FindIndex(func(r string, _ int) bool {
							return proxy[0] == r
						})
						if index > -1 {
							ctx.JSON(200, struct {
								Status  bool
								Message string
								Index   int
							}{
								Status: true,
								Index:  index,
							})
						} else {
							ctx.JSON(404, struct {
								Status  bool
								Message string
								Index   int
							}{
								Status:  false,
								Message: "not found",
								Index:   index,
							})
						}

					},
				},
				"locked": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "locked-proxy",
							Path:        "locked",
							Description: "locked proxies",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get stats of locked proxies")

						ctx.JSON(200, proxyCollection.locked)
					},
				},
				"active": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "proxy",
							Path:        "active",
							Description: "Active list",
						},
						Methods: []string{interfaces.GET},
						Middlewares: []interfaces.MiddlewareInterface{
							netHttp.UrlRequiredMiddleware,
						},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get active proxies")

						u := GetUrlParamGin(ctx)

						ctx.JSON(200, append(proxyCollection.locked[u.Hostname()], proxyCollection.allowed[u.Hostname()]...))
					},
				},
				"block": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "block",
							Path:        "block",
							Description: "block the proxy",
						},
						Methods: []string{interfaces.GET},
						Middlewares: []interfaces.MiddlewareInterface{
							netHttp.UrlRequiredMiddleware,
						},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						params := ctx.Request.URL.Query()
						blockProxy := params["proxy"]
						statusCode := params["code"]
						if len(blockProxy) == 0 || len(statusCode) == 0 {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: Errors.New("invalid request for /block").Error()})
							return
						}

						u := GetUrlParamGin(ctx)

						logger.Debug(fmt.Sprintf("block proxy by url %s, proxy: %s, status: %s", u.String(), blockProxy[0], statusCode[0]))

						err := proxyCollection.Block(blockProxy[0], u)
						if err != nil {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: err.Error()})
							return
						}

						ctx.JSON(200, struct {
							Success bool
						}{Success: true})
					},
				},
				"bad_request": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "bad",
							Path:        "bad_request",
							Description: "error code 500",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Some access errors.")
						ctx.JSON(500, nil)
					},
				},
				"stats": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "stats",
							Path:        "stats",
							Description: "Stats requests",
						},
						Methods: []string{interfaces.GET},
						Middlewares: []interfaces.MiddlewareInterface{
							netHttp.UrlRequiredMiddleware,
						},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Check stats.")

						u := GetUrlParamGin(ctx)

						payload := proxyCollection.getAvailablePayload(u)

						ctx.JSON(200, &Stats{
							Stats:           proxyCollection.stats,
							Active:          append(proxyCollection.locked[u.Hostname()], proxyCollection.allowed[u.Hostname()]...),
							Blocked:         proxyCollection.banned,
							Allowed:         proxyCollection.allowed[u.Hostname()],
							Locked:          proxyCollection.locked,
							List:            proxyCollection.list,
							TestEndpoint:    payload,
							ObservableHosts: proxyCollection.observableHosts,
						})
					},
				},
				"clear": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "clear",
							Path:        "clear",
							Description: "clear proxy history",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("clear proxy history.")
						e := proxyCollection.Clear()
						ctx.JSON(200, e)
					},
				},
				"check": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "check",
							Path:        "check",
							Description: "check domain by url",
						},
						Methods: []string{interfaces.POST},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("change url and headers for check available host through proxy")
						payload := &UpdateCheckPayload{}

						e := ctx.BindJSON(payload)
						if e != nil {
							ctx.JSON(500, e)
							return
						}
						_, u := payload.ParseUrl()

						proxyCollection.setAvailablePayload(payload)

						ctx.JSON(200, proxyCollection.getAvailablePayload(u))

					},
				},
			},
		},
	)
}
