package errors

import Errors "errors"

var ActionAlreadyExists  = Errors.New("action already exists")
var ServerNotFoundError = Errors.New("server not found error")
var ResourceAlreadyExist = Errors.New("resource already exist")
var NoResourcesAvailable  = Errors.New("no resources available")
var ServerMethodNotImplemented = Errors.New("server method not implemented")
var ServerOnStartError = Errors.New("server method on start not implemented")
var ActionContextRequestFailed  = Errors.New("action context request failed. ")
var ActionAddMethodNotImplemented = Errors.New("action add method not implemented")
var ActionPathNotFound  = Errors.New("action path not found in context. need set RoutePath")
var ServerOnHandlerMethodNotImplemented = Errors.New("server on handler method not implemented")

//var ResourceNotFound	 = Errors.New("resource not found")
//var ActionNotFound       = Errors.New("action not found")
//var httpParamsNotValid	 = Errors.New("http params isn't valid")
//var toketExpired		 = Errors.New("token expired")
//var sessionExpired		 = Errors.New("session expired")

var ServerEngineNotImplemented = Errors.New("server engine not implemented")
var ServerContextFailed = Errors.New("server context failed")
var TracerContextNotFound = Errors.New("tracer context not found. need set TracingOptions ")