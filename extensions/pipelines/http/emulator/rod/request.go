package rod

import (
	"bytes"
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/log"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	netHttp "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"github.com/vortex14/gotyphoon/interfaces"
)

func processElementsAfterPreLoad(
	logger interfaces.LoggerInterface,
	page *rod.Page,
	detailOptions *DetailsOptions) {

	if detailOptions != nil {
		if detailOptions.EventOptions.NetworkResponseReceived {
			detailOptions.EventOptions.Wait()
		}
	}

	if detailOptions != nil && detailOptions.SleepAfter > 0 {
		logger.Debug(fmt.Sprintf("Sleep after load: %f", detailOptions.SleepAfter))
		time.Sleep(time.Duration(detailOptions.SleepAfter) * time.Second)
	}
	if detailOptions != nil && len(detailOptions.MustElement) > 0 {

		if !detailOptions.Click && len(detailOptions.Input) == 0 {
			page.MustElement(detailOptions.MustElement)
		} else if detailOptions.Click {
			page.MustElement(detailOptions.MustElement).MustClick()
		} else if len(detailOptions.Input) > 0 {
			ie := page.MustElement(detailOptions.MustElement).Input(detailOptions.Input)
			if ie != nil {
				logger.Error(ie)
			}
		}
	}
}

func CreateRodRequestPipeline(
	opts *forms.Options,
	detailOptions *DetailsOptions) *pipelines.TaskPipeline {

	_middlewares := make([]interfaces.MiddlewareInterface, 0)

	if detailOptions.ProxyRequired {
		_middlewares = append(_middlewares, netHttp.ConstructorProxySettingMiddleware(true))
	}

	return &pipelines.TaskPipeline{
		BasePipeline: &forms.BasePipeline{
			NotIgnorePanic: false,
			MetaInfo: &label.MetaInfo{
				Name: "Rod http request",
			},
			Options:     opts,
			Middlewares: _middlewares,
		},

		Fn: func(context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface) (error, context.Context) {

			logger.Infof("RUN rod request proxy: %s , proxy_server: %s url: %s",
				task.GetProxyAddress(), task.GetProxyServerUrl(), task.GetFetcherUrl())

			detailOptions.Options.Timeout = task.GetFetcherTimeout()
			detailOptions.Options.Proxy = task.GetProxyAddress()

			if len(task.GetProxyAddress()) > 0 {
				context = log.PatchCtx(context, map[string]interface{}{"proxy": task.GetProxyAddress()})
				_, logger = log.Get(context)
			}

			browser := CreateBaseBrowser(detailOptions.Options)

			var pErr error

			errR := rod.Try(func() {

				browser = browser.MustConnect()

				browser.MustIgnoreCertErrors(true)

				page := browser.MustPage(task.GetFetcherUrl())

				processElementsAfterPreLoad(logger, page, detailOptions)

				logger.Debug("the page loaded")
				context = NewPageCtx(context, page)

				page.MustWaitLoad()
				body := page.MustHTML()

				doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(body)))
				if err != nil {
					pErr = fmt.Errorf("goquery: %s", err.Error())
				}

				context = html.NewHtmlCtx(context, doc)
				context = NewBodyResponse(context, &body)

				context = NewBrowserCtx(context, browser)

			})

			if errR != nil {
				pErr = errR

				if cE := rod.Try(func() {
					browser.MustClose()
				}); cE != nil {
					pErr = cE
				}

			}

			return pErr, context
		},
		Cn: func(err error,
			context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface) {

			if task.GetSaveData("SKIP_CN") == "skip" {
				return
			}

			if len(task.GetProxyAddress()) == 0 || len(task.GetProxyServerUrl()) == 0 {
				return
			}

			// Block current proxy
			if netHttp.MakeBlockRequest(logger, task) != nil {
				logger.Error("Fatal exception. Impossible to block the proxy.")
				os.Exit(1)
			}
		},
	}
}
