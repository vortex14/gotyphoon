package typhoon

import "context"

// Pipeline represents a blocking operation in a pipeline. Implementing `Pipeline` will allow you to add
// business logic to your pipelines without directly managing channels. This simplifies your unit tests
// and eliminates channel management related bugs.
type Pipeline interface {
	// Run processes an input and returns an output or an error, if the output could not be processed.
	// When the context is canceled, process should stop all blocking operations and return the `Context.Err()`.
	Run(ctx context.Context, i interface{}) (interface{}, error)

	// Cancel is called if process returns an error or if the context is canceled while there are still items in the `in <-chan interface{}`.
	Cancel(i interface{}, err error)
}

// NewPipeline creates a process and cancel func
func NewPipeline(
	process func(ctx context.Context, i interface{}) (interface{}, error),
	cancel func(i interface{}, err error),
) Pipeline {
	return &pipeline{process, cancel}
}

// pipeline implements Pipeline
type pipeline struct {
	process func(ctx context.Context, i interface{}) (interface{}, error)
	cancel  func(i interface{}, err error)
}

func (p *pipeline) Run(ctx context.Context, i interface{}) (interface{}, error) {
	return p.process(ctx, i)
}

func (p *pipeline) Cancel(i interface{}, err error) {
	p.cancel(i, err)
}
