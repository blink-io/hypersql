//go:build sqlite && cgo && sqlite_cgo

package hypersql

import (
	"database/sql/driver"

	"github.com/mattn/go-sqlite3"
)

func getRawSQLiteDriver() driver.Driver {
	return &sqlite3.SQLiteDriver{}
}
