package errors


import Errors "errors"

var BadRequest = Errors.New("bad request")
var ActionMethodsNotFound = Errors.New("action methods not found")