package timing

import (
	"context"
	"log"
	"time"

	"github.com/qustavo/sqlhooks/v2"
)

type ctxKey struct{}

type hook struct {
	logf func(string, ...any)
}

var _ sqlhooks.Hooks = (*hook)(nil)

func New(ops ...Option) sqlhooks.Hooks {
	h := new(hook)
	for _, o := range ops {
		o(h)
	}
	if h.logf == nil {
		h.logf = log.Printf
	}
	return h
}

func (h *hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	newCtx := context.WithValue(ctx, ctxKey{}, time.Now())
	return newCtx, nil
}

func (h *hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if before, ok := ctx.Value(ctxKey{}).(time.Time); ok {
		h.logf("[TIMING] Executed SQL, timing cost [%s] for: %s", time.Since(before), query)
	}
	return ctx, nil
}
