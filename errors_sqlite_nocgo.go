//go:build !use_cgo

package hypersql

import (
	"github.com/spf13/cast"
	"modernc.org/sqlite"
)

type SQLiteError = sqlite.Error

var sqliteErrorHandlers = map[int]func(*SQLiteError) *Error{
	// SQLITE_CONSTRAINT_CHECK (275)
	275: func(e *SQLiteError) *Error {
		return ErrConstraintCheck.
			Renew(cast.ToString(e.Code()), e.Error(), e)
	},
	// SQLITE_CONSTRAINT_FOREIGNKEY (787)
	787: func(e *SQLiteError) *Error {
		return ErrConstraintForeignKey.
			Renew(cast.ToString(e.Code()), e.Error(), e)
	},
	// SQLITE_CONSTRAINT_NOTNULL (1299)
	1299: func(e *SQLiteError) *Error {
		return ErrConstraintNotNull.
			Renew(cast.ToString(e.Code()), e.Error(), e)
	},
	// SQLITE_CONSTRAINT_PRIMARYKEY (1555).
	// Notes: In DBMS, primary key is a unique key too.
	1555: sqliteUniqueConstraintHandler,
	// SQLITE_CONSTRAINT_UNIQUE (2067)
	2067: sqliteUniqueConstraintHandler,
}

func RegisterSQLiteErrorHandler(number int, fn func(*SQLiteError) *Error) {
	sqliteErrorHandlers[number] = fn
}

func sqliteUniqueConstraintHandler(e *SQLiteError) *Error {
	return ErrConstraintUnique.
		Renew(cast.ToString(e.Code()), e.Error(), e)
}

func handleSQLiteError(e *SQLiteError) *Error {
	if h, ok := sqliteErrorHandlers[e.Code()]; ok {
		return h(e)
	} else {
		return ErrOther.Renew(cast.ToString(e.Code()), e.Error(), e)
	}
}
