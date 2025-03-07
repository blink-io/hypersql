package hypersql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
)

type MySQLError = mysql.MySQLError

var mysqlErrorHandlers = map[uint16]func(*mysql.MySQLError) *Error{
	// Error number: 1169; Symbol: ER_DUP_UNIQUE; SQLSTATE: 23000
	// Message: Can't write, because of unique constraint, to table '%s'
	1169: func(e *mysql.MySQLError) *Error {
		code := cast.ToString(e.Number)
		return ErrConstraintUnique.As(code, e.Message, e)
	},

	// Error number: 1172; Symbol: ER_TOO_MANY_ROWS; SQLSTATE: 42000
	// Message: Result consisted of more than one row
	1172: func(e *mysql.MySQLError) *Error {
		code := cast.ToString(e.Number)
		return ErrTooManyRows.As(code, e.Message, e)
	},

	// Error number: 1329; Symbol: ER_SP_FETCH_NO_DATA; SQLSTATE: 02000
	// Message: No data - zero rows fetched, selected, or processed
	1329: func(e *mysql.MySQLError) *Error {
		code := cast.ToString(e.Number)
		return ErrNoRows.As(code, e.Message, e)
	},

	// Error number: 1216; Symbol: ER_NO_REFERENCED_ROW; SQLSTATE: 23000
	// Message: Cannot add or update a child row: a foreign key constraint fails
	1216: mysqlFKConstraintErrHandler,

	// Error number: 1217; Symbol: ER_ROW_IS_REFERENCED; SQLSTATE: 23000
	// Message: Cannot delete or update a parent row: a foreign key constraint fails
	1217: mysqlFKConstraintErrHandler,

	// Error number: 1263; Symbol: ER_WARN_NULL_TO_NOTNULL; SQLSTATE: 22004
	1263: func(e *mysql.MySQLError) *Error {
		code := cast.ToString(e.Number)
		return ErrConstraintNotNull.As(code, e.Message, e)
	},

	// Error number: 1451; Symbol: ER_ROW_IS_REFERENCED_2; SQLSTATE: 23000
	// Message: Cannot delete or update a parent row: a foreign key constraint fails (%s)
	1451: mysqlFKConstraintErrHandler,

	// Error number: 1452; Symbol: ER_NO_REFERENCED_ROW_2; SQLSTATE: 23000
	// Message: Cannot add or update a child row: a foreign key constraint fails (%s)
	1452: mysqlFKConstraintErrHandler,

	// Error number: 3819; Symbol: ER_CHECK_CONSTRAINT_VIOLATED; SQLSTATE: HY000
	// Message: Check constraint '%s' is violated.
	// ER_CHECK_CONSTRAINT_VIOLATED was added in 8.0.16.
	3819: mysqlCheckConstraintErrHandler,

	// Error number: 3820; Symbol: ER_CHECK_CONSTRAINT_REFERS_UNKNOWN_COLUMN; SQLSTATE: HY000
	// Message: Check constraint '%s' refers to non-existing column '%s'.
	// ER_CHECK_CONSTRAINT_REFERS_UNKNOWN_COLUMN was added in 8.0.16.
	3820: mysqlCheckConstraintErrHandler,
}

func RegisterMySQLErrorHandler(number uint16, fn func(*mysql.MySQLError) *Error) {
	mysqlErrorHandlers[number] = fn
}

func mysqlFKConstraintErrHandler(e *mysql.MySQLError) *Error {
	code := cast.ToString(e.Number)
	return ErrConstraintForeignKey.As(code, e.Message, e)
}

func mysqlCheckConstraintErrHandler(e *mysql.MySQLError) *Error {
	code := cast.ToString(e.Number)
	return ErrConstraintForeignKey.As(code, e.Message, e)
}

// handleMySQLError transforms *mysql.MySQLError to *Error.
// Doc: https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
// About SQLState value: This value is a five-character string (for example, '42S02').
// SQLSTATE values are taken from ANSI SQL and ODBC and are more standardized than the numeric error codes.
// The first two characters of an SQLSTATE value indicate the error class:
//
// Class = '00' indicates success.
//
// Class = '01' indicates a warning.
//
// Class = '02' indicates "not found".
// This is relevant within the context of cursors,
// and is used to control what happens when a cursor reaches the end of a data set.
// This condition also occurs for SELECT ... INTO var_list statements that retrieve no rows.
//
// Class > '02' indicates an exception.
//
// For server-side errors, not all MySQL error numbers have corresponding SQLSTATE values.
// In these cases, 'HY000' (general error) is used.
// For client-side errors, the SQLSTATE value is always 'HY000' (general error),
// so it is not meaningful for distinguishing one client error from another.
func handleMySQLError(e *mysql.MySQLError) *Error {
	if h, ok := mysqlErrorHandlers[e.Number]; ok {
		return h(e)
	} else {
		return ErrOther.As(cast.ToString(e.Number), e.Message, e)
	}
}

func isMySQLError(e error) bool {
	_, ok := isTargetErr[*MySQLError](e)
	return ok
}
