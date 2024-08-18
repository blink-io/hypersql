//go:build use_cgo && !cgo_ext

package hypersql

import (
	"database/sql/driver"

	"github.com/mattn/go-sqlite3"
)

func getRawSQLiteDriver() driver.Driver {
	return &sqlite3.SQLiteDriver{}
}
