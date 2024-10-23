package slog

import (
	"log/slog"
	"testing"

	mssql "github.com/microsoft/go-mssqldb"
)

func TestLogger_1(t *testing.T) {
	logger := New(slog.Default(), slog.LevelInfo)

	logger.Printf("Hello, %s", "Heison")

	mssql.SetLogger(logger)
}
