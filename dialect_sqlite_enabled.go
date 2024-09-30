//go:build sqlite

package hypersql

import (
	"context"
	"database/sql/driver"
)

func GetSQLiteDSN(dialect string) (Dsner, error) {
	if !IsCompatibleSQLiteDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return func(ctx context.Context, c *Config) (string, error) {
		dsn := c.Host
		return dsn, nil
	}, nil
}

func IsCompatibleSQLiteDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatibleSQLiteDialects)
}

func GetSQLiteDriver(dialect string) (driver.Driver, error) {
	if IsCompatibleSQLiteDialect(dialect) {
		return getRawSQLiteDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetSQLiteConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	dsn := toSQLiteDSN(c)
	drv := wrapDriverHooks(getRawSQLiteDriver(), c.DriverHooks...)
	return &dsnConnector{dsn: dsn, driver: drv}, nil
}
