//go:build sqlite && cgo && sqlite_cgo

package hypersql

import (
	"github.com/mattn/go-sqlite3"
	"github.com/spf13/cast"
)

type SQLiteError = sqlite3.Error

var sqliteErrorHandlers = map[int]func(*SQLiteError) *Error{
	// SQLITE_CONSTRAINT_CHECK (275)
	275: func(e *SQLiteError) *Error {
		code := cast.ToString(int(e.ExtendedCode))
		return ErrConstraintCheck.As(code, e.Error(), e)
	},
	// SQLITE_CONSTRAINT_FOREIGNKEY (787)
	787: func(e *SQLiteError) *Error {
		code := cast.ToString(int(e.ExtendedCode))
		return ErrConstraintForeignKey.As(code, e.Error(), e)
	},
	// SQLITE_CONSTRAINT_NOTNULL (1299)
	1299: func(e *SQLiteError) *Error {
		code := cast.ToString(int(e.ExtendedCode))
		return ErrConstraintNotNull.As(code, e.Error(), e)
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
	code := cast.ToString(int(e.ExtendedCode))
	return ErrConstraintUnique.As(code, e.Error(), e)
}

// handleSQLiteError transforms *sqlite3.Error to *Error.
// Doc: https://www.sqlite.org/rescode.html
func handleSQLiteError(e *SQLiteError) *Error {
	code := int(e.ExtendedCode)
	if h, ok := sqliteErrorHandlers[code]; ok {
		return h(e)
	} else {
		return ErrOther.As(cast.ToString(code), e.Error(), e)
	}
}

func isSQLiteError(e error) bool {
	_, ok := isTargetErr[*SQLiteError](e)
	return ok
}
