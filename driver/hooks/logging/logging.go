package logging

import (
	"context"

	"github.com/qustavo/sqlhooks/v2"
)

type Func func(format string, args ...any)

var _ sqlhooks.Hooks = (Func)(nil)

func (f Func) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	f("Executed SQL, query: %s, args: %+v", query, args)
	return ctx, nil
}

func (f Func) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	return ctx, nil
}

type CtxFunc func(ctx context.Context, format string, args ...any)

var _ sqlhooks.Hooks = (CtxFunc)(nil)

func (f CtxFunc) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	f(ctx, "[LOGGING] Executed SQL, query: %s, args: %+v", query, args)
	return ctx, nil
}

func (f CtxFunc) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	return ctx, nil
}
