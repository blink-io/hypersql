//go:build sqlite && cgo && sqlite_cgo

package hypersql

import (
	"context"
	"database/sql/driver"

	"github.com/mattn/go-sqlite3"
)

var _ Validator = (*SQLiteExtra)(nil)

type SQLiteExtra struct {
}

func (c *SQLiteExtra) Validate(ctx context.Context) error {
	if c == nil {
		return ErrNilConfig
	}
	return nil
}

func RawSQLiteDriver() driver.Driver {
	return &sqlite3.SQLiteDriver{}
}
