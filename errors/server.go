package errors

import Errors "errors"

var ResourceAlreadyExist = Errors.New("resource already exist")
var NoResourcesAvailable  = Errors.New("no resources available")
var ActionAlreadyExists  = Errors.New("action already exists")

//var ResourceNotFound	 = Errors.New("resource not found")
//var ActionNotFound       = Errors.New("action not found")
//var httpParamsNotValid	 = Errors.New("http params isn't valid")
//var toketExpired		 = Errors.New("token expired")
//var sessionExpired		 = Errors.New("session expired")


