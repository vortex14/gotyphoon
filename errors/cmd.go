package errors

import Errors "errors"

var ErrorCmdNotFound = Errors.New("command for console line not found")
var ErrorStopCmd = Errors.New("command has error after finish")