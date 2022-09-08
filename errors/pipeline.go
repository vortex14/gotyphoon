package errors

import Errors "errors"

var ForceSkipPipelines = Errors.New("skip all next pipelines")
var PipelineContexFailed = Errors.New("invalid pipeline context")
var CtxLogFailed = Errors.New("context has not logger")
