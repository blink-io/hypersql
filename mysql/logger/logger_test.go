package logger

import (
	"fmt"
	"testing"
)

func TestLogger_1(t *testing.T) {
	ll := Logf(func(v ...any) {
		msg := fmt.Sprint(v...)
		fmt.Println(msg)
	})

	ll.Print(1, true, 12, "ok", 88, 3.14, 'C')
}
