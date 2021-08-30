package ctx

import (
	"context"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func UpdateContext(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, ContextKey(key), value)
}

func GetContextValue(ctx context.Context, key string) interface{} {
	return ctx.Value(ContextKey(key))
}
