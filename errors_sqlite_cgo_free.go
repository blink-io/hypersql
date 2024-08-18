//go:build no_cgo && cgo_free

package hypersql

import (
	"github.com/ncruces/go-sqlite3"
	"github.com/spf13/cast"
)

type SQLiteError = sqlite3.Error

var sqliteErrorHandlers = map[uint16]func(*SQLiteError) *Error{
	// SQLITE_CONSTRAINT_CHECK (275)
	275: func(e *SQLiteError) *Error {
		code := cast.ToString(uint16(e.ExtendedCode()))
		return ErrConstraintCheck.Renew(code, e.Error(), e)
	},
	// SQLITE_CONSTRAINT_FOREIGNKEY (787)
	787: func(e *SQLiteError) *Error {
		code := cast.ToString(uint16(e.ExtendedCode()))
		return ErrConstraintForeignKey.Renew(code, e.Error(), e)
	},
	// SQLITE_CONSTRAINT_NOTNULL (1299)
	1299: func(e *SQLiteError) *Error {
		code := cast.ToString(uint16(e.ExtendedCode()))
		return ErrConstraintNotNull.Renew(code, e.Error(), e)
	},
	// SQLITE_CONSTRAINT_PRIMARYKEY (1555).
	// Notes: In DBMS, primary key is a unique key too.
	1555: sqliteUniqueConstraintHandler,

	// SQLITE_CONSTRAINT_UNIQUE (2067)
	2067: sqliteUniqueConstraintHandler,
}

func RegisterSQLiteErrorHandler(number uint16, fn func(*SQLiteError) *Error) {
	sqliteErrorHandlers[number] = fn
}

func sqliteUniqueConstraintHandler(e *SQLiteError) *Error {
	code := cast.ToString(uint16(e.ExtendedCode()))
	return ErrConstraintUnique.Renew(code, e.Error(), e)
}

// handleSQLiteError transforms *sqlite3.Error to *Error.
// Doc: https://www.sqlite.org/rescode.html
func handleSQLiteError(e *SQLiteError) *Error {
	code := uint16(e.ExtendedCode())
	if h, ok := sqliteErrorHandlers[code]; ok {
		return h(e)
	} else {
		return ErrOther.Renew(cast.ToString(code), e.Error(), e)
	}
}
