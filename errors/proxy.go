package errors

import Errors "errors"

var ProxyBusy = Errors.New("proxy busy")
var ProxyBlocked = Errors.New("proxy blocked")
var ProxyTimeout = Errors.New("proxy timeout")
var ProxyAuthFailed = Errors.New("proxy auth failed")
var ProxyUrlWrong = Errors.New("proxy url wrong")
var ProxyTaskNotFound = Errors.New("task required proxy, but proxy not found in task")