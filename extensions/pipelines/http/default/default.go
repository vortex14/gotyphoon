package _default

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)

const (
	NAME        = "Default http Pipeline"
	DESCRIPTION = "Default logic for http pipeline"


	BASICAuthUserTaskKey = "login"
	BASICAuthPasswordTaskKey = "password"
)

// TODO: keep-alive connections

type HttpPipelineDefault struct {
	client          *http.Client
	request         *http.Request
	transport       *http.Transport

	response        *http.Response
	responseHeaders map[string]string

	*interfaces.BasePipeline
}

func (h *HttpPipelineDefault) setBasicAuth()  {
	if h.request != nil {
		h.request.SetBasicAuth(
			h.Task.Fetcher.Auth[BASICAuthUserTaskKey],
			h.Task.Fetcher.Auth[BASICAuthPasswordTaskKey],
		)
		//println("login: ", h.Task.Fetcher.Auth[BASICAuthUserTaskKey])
		//println("password: ", h.Task.Fetcher.Auth[BASICAuthPasswordTaskKey])
	}
	//curlTest := fmt.Sprintf("curl -u %s:%s %s",  h.Task.Fetcher.Auth[BASICAuthUserTaskKey],
	//	h.Task.Fetcher.Auth[BASICAuthPasswordTaskKey], h.Task.URL)
	//println(curlTest)
}

func (h *HttpPipelineDefault) setProxy()  {
	if h.Task.Fetcher.IsProxyRequired == true {

		proxyURL, err := url.Parse(h.Task.Fetcher.Proxy)

		if err != nil {
			fmt.Println(err)
		}

		if proxyURL.Host != "" && proxyURL.Port() != "" {
			h.transport.Proxy = http.ProxyURL(proxyURL)
		}

	}
}

func (h *HttpPipelineDefault) initTransport()  {
	if h.transport == nil {
		h.transport = &http.Transport{
			ResponseHeaderTimeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
			IdleConnTimeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
		}
	}

}

func (h *HttpPipelineDefault) initRequest()  {
	h.request, _ = http.NewRequest(h.Task.Fetcher.Method, h.Task.URL, nil)
	h.setHeader()
}

func (h *HttpPipelineDefault) updateResponseHeaders()  {
	responseHeaders := make(map[string]string)
	for key, value := range h.response.Header {
		count := len(value)
		isCount := count > 0
		isKey := len(key) > 0
		if isCount && isKey {
			responseHeaders[key] = strings.Join(value, "")
		}

	}

	h.Task.Fetcher.Response.Headers = responseHeaders
	h.Task.Fetcher.Response.Code = h.response.StatusCode
}

func (h *HttpPipelineDefault) initClient() {
	h.Context = context.Background()


	h.initTransport()
	h.initRequest()
	h.setBasicAuth()
	h.client = &http.Client{
		Transport: h.transport,
		Timeout: time.Duration(h.Task.Fetcher.Timeout) * time.Second,
	}
}

func (h *HttpPipelineDefault) setHeader()  {
	for key, element := range h.Task.Fetcher.Headers {
		h.request.Header.Add(
			key,
			element,
		)

	}
}

func (h *HttpPipelineDefault) getResponse() (error, []byte) {
	response, err := h.client.Do(h.request)
	if err != nil {
		h.LOG.Error("response Error ======= ! ", err)
		h.Task.Fetcher.Response.Code = 599
		return err, nil
	}

	h.response = response

	defer h.response.Body.Close()
	var reader io.ReadCloser
	switch h.response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(h.response.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer reader.Close()
	default:
		reader = h.response.Body
	}

	data, err := ioutil.ReadAll(reader)

	if err != nil {
		color.Red("%s", err.Error())
		return err, nil
	}


	h.Task.Fetcher.Response.Code = response.StatusCode

	return nil, data
}

func (h *HttpPipelineDefault) Run() (error, interface{}) {

	h.initClient()
	err, responseData := h.getResponse()

	//h.LOG.Debug(len(responseData), err, h.Task.Fetcher.Response.Code)

	return err, responseData
}

func (h *HttpPipelineDefault) Retry()  {

}

func (h *HttpPipelineDefault) Finish()  {

}

func Constructor(
	task *task.TyphoonTask,
	project interfaces.Project,

	) interfaces.BasePipelineInterface {

	return &HttpPipelineDefault{BasePipeline: &interfaces.BasePipeline{
		BasePipelineLabel: &interfaces.BasePipelineLabel{
			Name: NAME,
			Description: DESCRIPTION,
		},
		Task: task,
		Project: project,
		LOG: logrus.WithFields(logrus.Fields{
			"Pipeline": NAME,
		}),
	}}
}