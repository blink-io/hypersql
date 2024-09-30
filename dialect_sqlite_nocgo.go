//go:build sqlite && !sqlite_cgo

package hypersql

import (
	"database/sql/driver"

	"modernc.org/sqlite"
)

func getRawSQLiteDriver() driver.Driver {
	return &sqlite.Driver{}
}
