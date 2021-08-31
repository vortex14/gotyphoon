package errors

import Errors "errors"



var LambdaRequired = Errors.New("lambda is required for active base pipeline")

var TaskPipelineRequiredHandler = Errors.New("task pipeline is required handler")
var TaskPipelineRequiredCancelHandler = Errors.New("task pipeline is required cancel handler")