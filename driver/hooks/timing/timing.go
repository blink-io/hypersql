package timing

import (
	"context"
	"log"
	"time"
)

type ctxKey struct{}

type Hook struct {
	logf func(string, ...any)
}

func New(ops ...Option) *Hook {
	h := new(Hook)
	for _, o := range ops {
		o(h)
	}
	if h.logf == nil {
		h.logf = log.Printf
	}
	return h
}

func (h *Hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	newCtx := context.WithValue(ctx, ctxKey{}, time.Now())
	return newCtx, nil
}

func (h *Hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if before, ok := ctx.Value(ctxKey{}).(time.Time); ok {
		h.logf("[TIMING] Executed SQL, timing cost [%s] for: %s", time.Since(before), query)
	}
	return ctx, nil
}
