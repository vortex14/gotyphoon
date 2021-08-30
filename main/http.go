package main

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"
	netHttp "github.com/vortex14/gotyphoon/extensions/pipelines/http/middlewares/net-http"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
	"net/http"
)


const (
	NAME        = "Default http Pipeline"
	DESCRIPTION = "Default logic for http pipeline"
)

// TODO: keep-alive connections. Pass to Demon work

type HttpPipelineDefault struct {
	forms.BasePipeline
}


func (h *HttpPipelineDefault) getResponse(
	client *http.Client,
	request *http.Request) (
	error, []byte, *http.Response) {

	//response, err := client.Do(request)
	//if err != nil {
	//	h.LOG.Error("response Error ======= ! ", err)
	//	h.Task.Fetcher.Response.Code = 599
	//	return err, nil, nil
	//}
	//
	//defer response.Body.Close()
	//var reader io.ReadCloser
	//switch response.Header.Get("Content-Encoding") {
	//case "gzip":
	//	reader, err = gzip.NewReader(response.Body)
	//	if err != nil {
	//		h.LOG.Error(err.Error())
	//		return err, nil, nil
	//	}
	//	defer reader.Close()
	//default:
	//	reader = response.Body
	//}
	//
	//data, err := ioutil.ReadAll(reader)
	//
	//if err != nil {
	//	logrus.Error(err.Error())
	//	return err, nil, nil
	//}



	return nil, nil, nil
}

func (h *HttpPipelineDefault) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	//transport := &http.Transport{
	//	ResponseHeaderTimeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
	//	IdleConnTimeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
	//}
	//
	//client := &http.Client{
	//	Transport: transport,
	//	Timeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
	//}
	//
	//request, err := http.NewRequest(h.Task.Fetcher.Method, h.Task.URL, nil)
	//
	//if err != nil {
	//	return err, nil
	//}
	//
	//
	//// Before http request. First step for prefetching
	//ctx := context.WithValue(h.Context, interfaces.ContextKey(netHttp.TASK), h.Task)
	//ctx = context.WithValue(ctx, interfaces.ContextKey(netHttp.CLIENT), client)
	//ctx = context.WithValue(ctx, interfaces.ContextKey(netHttp.REQUEST), request)
	//ctx = context.WithValue(ctx, interfaces.ContextKey(netHttp.TRANSPORT), transport)
	//
	//var middlewareException error
	//h.RunMiddlewareStack(ctx, func(middleware interfaces.MiddlewareInterface, err error) {
	//	h.LOG.Error(err.Error())
	//	middlewareException = err
	//})
	//
	//
	//err, data, response := h.getResponse(client, request)
	//println(err, string(data), response.StatusCode)
	////After http request. Define response
	////u := utils.Utils{}
	////println(111,u.PrintPrettyJson(h.Task.Fetcher.Headers))
	//
	//
}

func Constructor(
	task *task.TyphoonTask,
	project interfaces.Project,

) interfaces.BasePipelineInterface {

	return &HttpPipelineDefault{
		BasePipeline: forms.BasePipeline{
			BaseLabel: interfaces.BaseLabel{
				Name:        NAME,
				Description: DESCRIPTION,
			},
			Middlewares: []interfaces.MiddlewareInterface{

				//netHttp.ConstructorProxyMiddleware(false),
				netHttp.ConstructorBasicAuthMiddleware(false),
				netHttp.ConstructorRequestHeaderMiddleware(true),
			},
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {


				//transport := &http.Transport{
				//	ResponseHeaderTimeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
				//	IdleConnTimeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
				//}
				//
				//client := &http.Client{
				//	Transport: transport,
				//	Timeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
				//}
				//
				//request, err := http.NewRequest(h.Task.Fetcher.Method, h.Task.URL, nil)

				return nil, nil
			},
		},
	}
}
