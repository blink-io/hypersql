package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/qustavo/sqlhooks/v2"
)

type hook struct {
	hub *sentry.Hub
}

var _ sqlhooks.Hooks = (*hook)(nil)
var _ sqlhooks.OnErrorer = (*hook)(nil)

func New(ops ...Option) (sqlhooks.Hooks, error) {
	h := new(hook)
	for _, o := range ops {
		o(h)
	}
	return h, nil
}

func (h *hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	if err != nil {
		h.hub.CaptureException(err)
	}
	return err
}

func (h *hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return ctx, nil
}

func (h *hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return ctx, nil
}
