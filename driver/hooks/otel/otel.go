package otel

import (
	"context"
)

type Hook struct {
}

func New(ops ...Option) *Hook {
	h := &Hook{}
	return h
}

func (h *Hook) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return ctx, nil
}

func (h *Hook) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	return ctx, nil
}
