package errors

import Errors "errors"

var BadRequest = Errors.New("bad request")
var ActionMethodsNotFound = Errors.New("action method or pipeline not found")