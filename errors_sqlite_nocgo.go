//go:build sqlite && !sqlite_cgo

package hypersql

import (
	"github.com/spf13/cast"
	sqlite3 "modernc.org/sqlite"
)

type SQLiteError = sqlite3.Error

var sqliteErrorHandlers = map[int]func(*SQLiteError) *Error{
	// SQLITE_CONSTRAINT_CHECK (275)
	275: func(e *SQLiteError) *Error {
		return ErrConstraintCheck.
			As(cast.ToString(e.Code()), e.Error(), e)
	},
	// SQLITE_CONSTRAINT_FOREIGNKEY (787)
	787: func(e *SQLiteError) *Error {
		return ErrConstraintForeignKey.
			As(cast.ToString(e.Code()), e.Error(), e)
	},
	// SQLITE_CONSTRAINT_NOTNULL (1299)
	1299: func(e *SQLiteError) *Error {
		return ErrConstraintNotNull.
			As(cast.ToString(e.Code()), e.Error(), e)
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
		As(cast.ToString(e.Code()), e.Error(), e)
}

func handleSQLiteError(e *SQLiteError) *Error {
	if h, ok := sqliteErrorHandlers[e.Code()]; ok {
		return h(e)
	} else {
		return ErrOther.As(cast.ToString(e.Code()), e.Error(), e)
	}
}

func isSQLiteError(e error) bool {
	_, ok := isTargetErr[*SQLiteError](e)
	return ok
}
