//go:build use_cgo && cgo_ext

package hypersql

import (
	"database/sql/driver"

	sqlite3drv "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func getRawSQLiteDriver() driver.Driver {
	return &sqlite3drv.SQLite{}
}
