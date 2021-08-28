package forms

import (
	"github.com/vortex14/gotyphoon/interfaces"
	"time"
)

type PipelineGroup struct {
	*interfaces.BaseLabel

	errorCount    int64
	duration      time.Time
	timeLife      time.Time

	LambdaMap     map[string]interfaces.LambdaInterface
	PyLambdaMap   map[string]interfaces.LambdaInterface

	Stages        []interfaces.BasePipelineInterface
	Consumers     map[string]interfaces.ConsumerInterface

}
