package errors

import (
	Errors "errors"

	"github.com/pkg/errors"
)

var (
	PipelineCrowded      = errors.New("pipeline crowded")
	ForceSkipPipelines   = Errors.New("skip all next pipelines")
	PipelineContexFailed = Errors.New("invalid pipeline context")
	CtxLogFailed         = Errors.New("context has not logger")
)