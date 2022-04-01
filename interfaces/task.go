package interfaces

type TaskInterface interface {
	IsMaxRetry() bool
	UpdateRetriesCounter()
	IsRetry() bool
	TaskFetcherInterface
}
type TaskFetcherInterface interface {
	GetFetcherMethod() string
	GetFetcherTimeout() int

	GetFetcherUrl() string
	SetFetcherUrl(url string)

	SetStatusCode(code int)
	SetProxyAddress(address string)
	SetProxyServerUrl(url string)

	GetProxyAddress() string
	IsProxyRequired() bool

	SetUserAgent(agent string)
	GetUserAgent() string
}

type TaskProcessorInterface interface {
}

type TaskResultTransporterInterface interface {
}

type TaskSchedulerInterface interface {
}
