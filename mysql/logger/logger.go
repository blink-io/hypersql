package logger

import "github.com/go-sql-driver/mysql"

type Logf func(v ...any)

func (l Logf) Print(v ...any) {
	l(v...)
}

var _ mysql.Logger = (Logf)(nil)
