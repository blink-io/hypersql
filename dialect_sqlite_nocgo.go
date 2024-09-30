//go:build sqlite && !sqlite_cgo

package hypersql

import (
	"database/sql/driver"

	sqlite3 "modernc.org/sqlite"
)

func getRawSQLiteDriver() driver.Driver {
	return &sqlite3.Driver{}
}
