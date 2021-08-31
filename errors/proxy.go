package errors

import Errors "errors"

var ProxyBusy = Errors.New("proxy busy")
var ProxyBlocked = Errors.New("proxy blocked")
var ProxyTimeout = Errors.New("proxy timeout")
var ProxyAuthFailed = Errors.New("proxy auth failed")
var ProxyUrlWrong = Errors.New("proxy url wrong")
var ProxyTaskNotFound = Errors.New("task required proxy, but proxy not found in task")
var ProxyTaskRequired = Errors.New("proxy required for task. task.Fetcher.IsProxyRequired = false")
var ProxyServerNotAvailable = Errors.New("proxy server isn't available")
var ProxyServerError        = Errors.New("response from proxy server isn't correct. ")