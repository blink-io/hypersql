package slog

import (
	"context"
	"fmt"
	"log/slog"

	mssql "github.com/microsoft/go-mssqldb"
)

var _ mssql.Logger = (*Logger)(nil)

type Logger struct {
	sl  *slog.Logger
	lv  slog.Level
	ctx context.Context
}

func New(sl *slog.Logger, lv slog.Level) *Logger {
	l := &Logger{
		sl:  sl,
		lv:  lv,
		ctx: context.Background(),
	}
	return l
}

func (l *Logger) Printf(format string, args ...any) {
	l.sl.Log(l.ctx, l.lv, fmt.Sprintf(format, args...))
}

func (l *Logger) Println(v ...any) {
	l.sl.Log(l.ctx, l.lv, fmt.Sprint(v...))
}
