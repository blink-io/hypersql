package hypersql

import (
	"context"
	"database/sql/driver"

	"github.com/marcboeker/go-duckdb"
)

var compatibleDuckDBDialects = []string{
	DialectDuckDB,
}

func init() {
	d := DialectDuckDB
	connectors[d] = GetDuckDBConnector
}

type DuckDBConfig struct {
	Calendar string
	TimeZone string
}

func (c *DuckDBConfig) Validate(ctx context.Context) error {
	if c == nil {

	}
	return nil
}

func GetDuckDBDSN(dialect string) (Dsner, error) {
	if !IsCompatibleDuckDBDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return func(ctx context.Context, c *Config) (string, error) {
		dsn := toDuckDBDSN(c)
		return dsn, nil
	}, nil
}

func GetDuckDBDriver(dialect string) (driver.Driver, error) {
	if IsCompatibleDuckDBDialect(dialect) {
		return getRawDuckDBDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetDuckDBConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	dsn := toDuckDBDSN(c)
	drv := wrapDriverHooks(getRawDuckDBDriver(), c.DriverHooks...)
	return &dsnConnector{dsn: dsn, driver: drv}, nil
}

func IsCompatibleDuckDBDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatibleDuckDBDialects)
}

func getRawDuckDBDriver() driver.Driver {
	return &duckdb.Driver{}
}

func toDuckDBDSN(c *Config) string {
	dsn := c.Host
	return dsn
}
