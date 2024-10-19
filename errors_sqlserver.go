package hypersql

import (
	mssql "github.com/microsoft/go-mssqldb"
)

type SQLServerError = mssql.Error

func isSQLServerError(e error) bool {
	_, ok1 := isTargetErr[*SQLServerError](e)
	_, ok2 := isTargetErr[SQLServerError](e)
	return ok1 || ok2
}
