package slog

import (
	"log/slog"
	"testing"
)

func TestLogger_1(t *testing.T) {
	logger := New(slog.Default(), slog.LevelInfo)

	logger.Printf("Hello, %s", "Heison")
}
