package net_http

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"

	"context"
	"net/http"
	"testing"
)

func init() {
	log.InitD()
}

func TestRetryHttpPipeline(t *testing.T) {

	l := log.New(map[string]interface{}{"test": "test"})
	ctxl := log.NewCtx(context.Background(), l)

	Convey("create a http pipeline and request to unavailable host", t, func() {
		countIter := 0

		_task := fake.CreateDefaultTask()

		_task.SetFetcherUrl("https://2ip2.ru")

		preparedCtx := task.PatchCtx(ctxl, _task)

		var errP error
		p1 := &HttpRequestPipeline{
			BasePipeline: &forms.BasePipeline{
				Middlewares: []interfaces.MiddlewareInterface{
					ConstructorPrepareRequestMiddleware(true),
				},
				NotIgnorePanic: true,
				MetaInfo: &label.MetaInfo{
					Name: "http-request",
				},
			},

			Fn: func(context context.Context,
				task interfaces.TaskInterface,
				logger interfaces.LoggerInterface,
				client *http.Client,
				request *http.Request,
				transport *http.Transport) (error, context.Context) {

				countIter += 1

				err, response, data := Request(client, request, logger)
				if err != nil {
					return err, nil
				}
				context = NewResponseCtx(context, response)
				context = NewResponseDataCtx(context, data)

				return nil, context

			},
			Cn: func(err error,
				context context.Context,
				task interfaces.TaskInterface,
				logger interfaces.LoggerInterface) {

				logger.Error("--- ", err.Error())
			},
		}

		var errM error
		p1.RunMiddlewareStack(preparedCtx, func(middleware interfaces.MiddlewareInterface, err error) {
			errM = err
			p1.Cancel(preparedCtx, l, err)
		}, func(ctx context.Context) {
			preparedCtx = ctx
		})

		So(errM, ShouldBeNil)

		p1.Run(preparedCtx, func(pipeline interfaces.BasePipelineInterface, err error) {
			errP = err
		}, func(ctx context.Context) {
			preparedCtx = ctx
		})

		l.Debug(errP)

		So(errP, ShouldBeError)

		status, _, data := GetResponseCtx(preparedCtx)

		So(status, ShouldBeFalse)

		So(data, ShouldBeNil)

		So(countIter, ShouldEqual, 7)
	})

	Convey("create a http pipeline and request to unavailable host at once", t, func() {
		countIter := 0

		_task := fake.CreateDefaultTask()

		_task.SetFetcherUrl("https://2ip2.ru")

		ctxt := task.PatchCtx(ctxl, _task)

		var errP error
		var preparedCtx context.Context

		p1 := &HttpRequestPipeline{
			BasePipeline: &forms.BasePipeline{
				Middlewares: []interfaces.MiddlewareInterface{
					ConstructorPrepareRequestMiddleware(true),
				},
				NotIgnorePanic: true,
				Options:        forms.GetNotRetribleOptions(),
				MetaInfo: &label.MetaInfo{
					Name: "http-request",
				},
			},

			Fn: func(context context.Context,
				task interfaces.TaskInterface,
				logger interfaces.LoggerInterface,
				client *http.Client,
				request *http.Request,
				transport *http.Transport) (error, context.Context) {

				countIter += 1

				err, response, data := Request(client, request, logger)
				if err != nil {
					return err, nil
				}
				context = NewResponseCtx(context, response)
				context = NewResponseDataCtx(context, data)

				return nil, context

			},
			Cn: func(err error,
				context context.Context,
				task interfaces.TaskInterface,
				logger interfaces.LoggerInterface) {

				logger.Error("--- ", err.Error())
			},
		}

		var errM error
		p1.RunMiddlewareStack(ctxt, func(middleware interfaces.MiddlewareInterface, err error) {
			errM = err
			p1.Cancel(ctxt, l, err)
		}, func(ctx context.Context) {
			preparedCtx = ctx
		})

		So(errM, ShouldBeNil)

		p1.Run(preparedCtx, func(pipeline interfaces.BasePipelineInterface, err error) {
			errP = err
		}, func(ctx context.Context) {
			preparedCtx = ctx
		})

		l.Debug(errP)

		So(errP, ShouldBeError)

		status, _, data := GetResponseCtx(preparedCtx)

		So(status, ShouldBeFalse)

		So(data, ShouldBeNil)

		So(countIter, ShouldEqual, 1)
	})

	Convey("create a http pipeline with unavailable proxy", t, func() {
		countIter := 0

		_task := fake.CreateDefaultTask()

		_task.SetFetcherUrl("https://2ip.ru")

		_task.SetProxyAddress("http://localhost:1414")
		_task.SetProxyServerUrl("")

		ctxt := task.PatchCtx(ctxl, _task)

		var errP error
		var preparedCtx context.Context

		So(errP, ShouldBeNil)

		p1 := &HttpRequestPipeline{
			BasePipeline: &forms.BasePipeline{

				NotIgnorePanic: true,
				Options:        forms.GetNotRetribleOptions(),
				MetaInfo: &label.MetaInfo{
					Name: "http-request",
				},
				Middlewares: []interfaces.MiddlewareInterface{
					ConstructorPrepareRequestMiddleware(true),
					ConstructorProxyRequestSettingsMiddleware(true),
				},
			},

			Fn: func(context context.Context,
				task interfaces.TaskInterface,
				logger interfaces.LoggerInterface,
				client *http.Client,
				request *http.Request,
				transport *http.Transport) (error, context.Context) {

				countIter += 1

				err, response, data := Request(client, request, logger)
				if err != nil {
					return err, nil
				}
				context = NewResponseCtx(context, response)
				context = NewResponseDataCtx(context, data)

				return nil, context

			},
			Cn: func(err error,
				context context.Context,
				task interfaces.TaskInterface,
				logger interfaces.LoggerInterface) {

				logger.Error("--- ", err.Error())
			},
		}
		var errM error
		p1.RunMiddlewareStack(ctxt, func(middleware interfaces.MiddlewareInterface, err error) {
			errM = err
			p1.Cancel(ctxt, l, err)
		}, func(ctx context.Context) {
			preparedCtx = ctx
		})

		So(errM, ShouldBeNil)

		var errR error
		p1.Run(preparedCtx, func(pipeline interfaces.BasePipelineInterface, err error) {
			errR = err
		}, func(ctx context.Context) {
			preparedCtx = ctx
		})

		l.Debug(errR)
		status, _, _ := GetResponseCtx(preparedCtx)

		So(status, ShouldBeFalse)
		So(errR, ShouldBeError)

		So(countIter, ShouldEqual, 1)
	})
}
