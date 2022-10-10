package rod

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"os"
	"time"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"
)

func CreateProxyRodRequestPipeline(opts *forms.Options, evopts *EventOptions) *HttpRodRequestPipeline {

	return &HttpRodRequestPipeline{
		BasePipeline: &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name: "Rod http request",
			},
			Options: opts,
			Middlewares: []interfaces.MiddlewareInterface{
				net_http.ConstructorProxySettingMiddleware(true),
				ConstructorRodProxyRequestMiddleware(true),
			},
		},

		Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface,
			browser *rod.Browser) (error, context.Context) {

			logger.Info(fmt.Sprintf("RUN rod request proxy: %s , proxy_server: %s url: %s", task.GetProxyAddress(), task.GetProxyServerUrl(), task.GetFetcherUrl()))

			page := browser.DefaultDevice(devices.IPhoneX).
				Timeout(time.Duration(task.GetFetcherTimeout()) * time.Second).
				MustConnect().
				MustPage(task.GetFetcherUrl())

			if evopts != nil {
				evopts.Wait()
			}

			logger.Debug("page opened")
			page.MustWaitLoad()
			logger.Debug("the page loaded")
			context = NewPageCtx(context, page)

			body := page.MustHTML()

			doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(body)))
			if err != nil {
				return err, context
			}

			context = html.NewHtmlCtx(context, doc)
			context = NewBodyResponse(context, &body)

			page.MustClose()

			return nil, context
		},
		Cn: func(err error,
			context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface) {

			if task.GetSaveData("SKIP_CN") == "skip" {
				return
			}

			// Block current proxy
			if net_http.MakeBlockRequest(logger, task) != nil {
				logger.Error("Fatal exception. Impossible to block the proxy.")
				os.Exit(1)
			}
		},
	}
}
