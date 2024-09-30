//go:build !sqlite

package hypersql

import (
	"context"
	"database/sql/driver"

	"github.com/spf13/cast"
)

type SQLiteError struct {
	Code int
}

func (e *SQLiteError) Error() string {
	return ""
}

func GetSQLiteDSN(dialect string) (Dsner, error) {
	return nil, ErrUnsupportedDialect
}

func IsCompatibleSQLiteDialect(dialect string) bool {
	return false
}

func GetSQLiteDriver(dialect string) (driver.Driver, error) {
	return nil, ErrUnsupportedDriver
}

func GetSQLiteConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	return nil, ErrUnsupportedDriver
}

func handleSQLiteError(e *SQLiteError) *Error {
	return ErrOther.Renew(cast.ToString(e.Code), e.Error(), e)
}
