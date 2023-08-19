package errors

import Errors "errors"

var (
	ActionAlreadyExists  = Errors.New("action already exists")
	ServerNotFoundError  = Errors.New("server not found error")
	ResourceAlreadyExist = Errors.New("resource already exist")

	NoResourcesAvailable       = Errors.New("no resources available")
	ServerMethodNotImplemented = Errors.New("server method not implemented")
	ServerOnStartError         = Errors.New("server method on start not implemented")
	ActionContextRequestFailed = Errors.New("action context request failed. ")

	ActionAddMethodNotImplemented = Errors.New("action add method not implemented")

	ActionPathNotFound    = Errors.New("action path not found in context. need set RoutePath")
	ActionFailed          = Errors.New("action failed")
	ActionErrRequestModel = Errors.New("request model isn't correct for action")

	ServerOnHandlerMethodNotImplemented = Errors.New("server on handler method not implemented")

	ServerEngineNotImplemented = Errors.New("server engine not implemented")

	ServerContextFailed = Errors.New("server context failed")

	TracerContextNotFound = Errors.New("tracer context not found. need set TracingOptions ")
)

//var ResourceNotFound	 = Errors.New("resource not found")
//var ActionNotFound       = Errors.New("action not found")
//var httpParamsNotValid	 = Errors.New("http params isn't valid")
//var toketExpired		 = Errors.New("token expired")
//var sessionExpired		 = Errors.New("session expired")
