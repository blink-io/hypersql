package hypersql

import (
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/spf13/cast"
)

type SQLServerError = mssql.Error

var sqlServerErrorHandlers = map[int32]func(*mssql.Error) *Error{
	2627: func(e *mssql.Error) *Error {
		return ErrOther.As(cast.ToString(e.Number), e.Message, e)
	},
}

func handleSQLServerError(e *mssql.Error) *Error {
	if h, ok := sqlServerErrorHandlers[e.Number]; ok {
		return h(e)
	} else {
		return ErrOther.As(cast.ToString(e.Number), e.Message, e)
	}
}

func isSQLServerError(e error) bool {
	_, ok1 := isTargetErr[*SQLServerError](e)
	_, ok2 := isTargetErr[SQLServerError](e)
	return ok1 || ok2
}
