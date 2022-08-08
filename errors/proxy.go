package errors

import Errors "errors"

var (
	ProxyBusy               = Errors.New("proxy busy")
	ProxyBlocked            = Errors.New("proxy blocked")
	ProxyTimeout            = Errors.New("proxy timeout")
	ProxyAuthFailed         = Errors.New("proxy auth failed")
	ProxyUrlWrong           = Errors.New("proxy url wrong")
	ProxyTaskNotFound       = Errors.New("task required proxy, but proxy not found in task")
	ProxyTaskRequired       = Errors.New("proxy required for task. task.Fetcher.IsProxyRequired = false")
	ProxyServerNotAvailable = Errors.New("proxy server isn't available")
	ProxyServerError        = Errors.New("response from proxy server isn't correct. ")
	ProxyServerNotFound     = Errors.New("not found proxy server url")
)
