package errors

import Errors "errors"

var PipelineContexFailed = Errors.New("invalid pipeline context")
var CtxLogFailed = Errors.New("context has not logger")
