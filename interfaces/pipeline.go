package interfaces

import (
	"context"
)


type BasePipelineInterface interface {
	Run(ctx context.Context) (error, context.Context)
	RunMiddlewareStack(
		context context.Context,
		reject func(middleware MiddlewareInterface,err error),
		next func(ctx context.Context),
	)
	GetName() string
	GetDescription()string
}

type ProcessorPipelineInterface interface {
	BasePipelineInterface
	Crawl()
	Finish()
	Switch()
}

type BaseLabel struct {
	Name        string
	Description string
	Required    bool
}

func (p *BaseLabel) IsRequired() bool {
	return p.Required
}

func (p *BaseLabel) GetName() string {
	return p.Name
}

func (p *BaseLabel) GetDescription() string {
	return p.Description
}


type CallbackPipelineInterface interface {
	Call(ctx context.Context, data interface{})
}



type ConsumerInterface interface {

}

type LambdaInterface interface {

}

type HandlerInterface interface {
	
}

type ResponseInterface interface {
	
}