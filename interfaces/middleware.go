package interfaces

import (
	"context"
)

type ConstructorMiddleware func (required bool) MiddlewareInterface

type MiddlewareCallback func(
		context context.Context,
		logger LoggerInterface,
		reject func(err error),
		next func(ctx context.Context),
	)

type MiddlewareInterface interface {
	IsRequired() bool
	Pass(context context.Context,
		logger LoggerInterface,
		reject func(err error),
		next func(context context.Context),
	)

	MetaDataInterface
}



