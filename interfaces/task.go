package interfaces

import (
	"bytes"
	"net/url"
	"time"
)

type TaskInterface interface {
	IsMaxRetry() bool
	UpdateRetriesCounter()
	String() string
	IsRetry() bool
	TaskFetcherInterface
}
type TaskFetcherInterface interface {
	GetFetcherMethod() string
	SetFetcherMethod(method string)

	AddHeader(key string, value string)
	SetHeaders(headers map[string]string)

	GetFetcherTimeout() time.Duration
	SetFetcherTimeout(seconds int)

	GetFetcherUrl() string
	GetBase64FetcherURL() string
	SetFetcherUrl(url string)

	SetCookies(cookies string)
	GetCookies() string

	SetStatusCode(code int)
	SetProxyAddress(address string)
	SetProxyServerUrl(url string)
	GetProxyServerUrl() string

	GetProxyAddress() string
	IsProxyRequired() bool

	SetUserAgent(agent string)
	GetUserAgent() string

	SetJsonRequestData(values interface{})
	SetRequestBody(values url.Values)
	GetRequestBody() *bytes.Buffer

	SetSaveData(key string, value interface{})
	GetSaveData(key string) interface{}
}

type TaskProcessorInterface interface {
	SetSaveData(key string, value string)
	GetSaveData(key string) string
}

type TaskResultTransporterInterface interface {
}

type TaskSchedulerInterface interface {
}
