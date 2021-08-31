package errors

import Errors "errors"

var MiddlewareBasicAuthOptionsNotFound = Errors.New("middleware basic auth options not found")
var MiddlewareNotImplemented = Errors.New("middleware not implemented")
var MiddlewareRequired = Errors.New("middleware required")
var ForceSkipMiddlewares = Errors.New("skip all future middleware stack")

var MiddlewareContextFailed = Errors.New("invalid middleware context. ctx must contain task & request")