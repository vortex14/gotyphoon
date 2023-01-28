package rod

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"github.com/vortex14/gotyphoon/interfaces"
)

func getDevice(detailOptions *DetailsOptions) devices.Device {
	var device devices.Device
	if detailOptions.Device == nil {
		device = devices.IPadPro
	} else {
		device = *detailOptions.Device
	}
	return device
}

func createPageFromTask(browser *rod.Browser, task interfaces.TaskInterface, detailOptions *DetailsOptions) *rod.Page {
	return browser.DefaultDevice(getDevice(detailOptions)).
		Timeout(time.Duration(task.GetFetcherTimeout()) * time.Second).
		MustConnect().
		MustPage(task.GetFetcherUrl())
}

func processElementsAfterPreLoad(logger interfaces.LoggerInterface, page *rod.Page, detailOptions *DetailsOptions) {

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

func CreateProxyRodRequestPipeline(opts *forms.Options, detailOptions *DetailsOptions) *HttpRodRequestPipeline {

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

			page := createPageFromTask(browser, task, detailOptions)

			processElementsAfterPreLoad(logger, page, detailOptions)

			logger.Debug("the page loaded")
			context = NewPageCtx(context, page)

			page.MustWaitLoad()
			body := page.MustHTML()

			doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(body)))
			if err != nil {
				return fmt.Errorf("goquery: %s", err.Error()), context
			}

			context = html.NewHtmlCtx(context, doc)
			context = NewBodyResponse(context, &body)

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

func CreateRodRequestPipeline(opts *forms.Options, detailOptions *DetailsOptions) *HttpRodRequestPipeline {

	return &HttpRodRequestPipeline{
		BasePipeline: &forms.BasePipeline{
			NotIgnorePanic: true,
			MetaInfo: &label.MetaInfo{
				Name: "Rod http request",
			},
			Options: opts,
			Middlewares: []interfaces.MiddlewareInterface{
				ConstructorRodBasicRequestMiddleware(true),
			},
		},

		Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface,
			browser *rod.Browser) (error, context.Context) {

			logger.Info(fmt.Sprintf("RUN rod request to url: %s", task.GetFetcherUrl()))

			page := createPageFromTask(browser, task, detailOptions)

			logger.Debug("created a new page")

			processElementsAfterPreLoad(logger, page, detailOptions)

			context = NewPageCtx(context, page)

			page.MustWaitLoad()
			logger.Debug("the page loaded")
			body := page.MustHTML()

			doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(body)))
			if err != nil {
				return fmt.Errorf("goquery: %s", err.Error()), context
			}

			context = html.NewHtmlCtx(context, doc)
			context = NewBodyResponse(context, &body)

			return nil, context
		},
		Cn: func(err error,
			context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface) {

			if task.GetSaveData("SKIP_CN") == "skip" {
				return
			}
		},
	}
}
