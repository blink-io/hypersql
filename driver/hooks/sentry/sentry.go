package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/qustavo/sqlhooks/v2"
)

type hook struct {
	hub *sentry.Hub
}

var _ interface {
	sqlhooks.Hooks
	sqlhooks.OnErrorer
} = (*hook)(nil)

func New(ops ...Option) (sqlhooks.Hooks, error) {
	h := new(hook)
	for _, o := range ops {
		o(h)
	}
	if h.hub == nil {
		h.hub = sentry.CurrentHub()
	}
	return h, nil
}

func (h *hook) OnError(ctx context.Context, err error, query string, args ...any) error {
	if err != nil {
		h.hub.CaptureException(err)
	}
	return err
}

func (h *hook) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return ctx, nil
}

func (h *hook) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	return ctx, nil
}
