package hypersql

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/tracelog"
)

func TestLogger_1(t *testing.T) {
	var l Logger = func(format string, args ...any) {
		s := fmt.Sprintf(format, args...)
		fmt.Println(s)
	}
	lf := doLoggerFunc(l)

	lf(context.Background(), tracelog.LogLevelInfo, "Hello", map[string]interface{}{
		"level":   "INFO",
		"enabled": true,
		"point":   4.44,
		"score":   88,
	})
}
