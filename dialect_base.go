package hypersql

import (
	"context"
	"database/sql/driver"
	"slices"
	"strings"

	"github.com/qustavo/sqlhooks/v2"
)

type (
	Dsner = func(context.Context, *Config) (string, error)

	ConnectorFunc func(ctx context.Context, c *Config) (driver.Connector, error)

	//GetDriverFunc func(dialect string) (driver.Driver, error)

	//GetDSNFunc func(dialect string) (Dsner, error)
)

var (
	//drivers = make(map[string]GetDriverFunc)

	//dsners = make(map[string]GetDSNFunc)

	connectors = make(map[string]ConnectorFunc)
)

func GetFormalDialect(dialect string) string {
	if d, ok := IsCompatibleDialect(dialect); ok {
		return d
	}
	return ""
}

// IsCompatibleDialect checks
func IsCompatibleDialect(dialect string) (string, bool) {
	if IsCompatiblePostgresDialect(dialect) {
		return DialectPostgres, true
	} else if IsCompatibleMySQLDialect(dialect) {
		return DialectMySQL, true
	} else if IsCompatibleSQLiteDialect(dialect) {
		return DialectSQLite, true
	}
	return "", false
}

func isCompatibleDialectIn(dialect string, compatibleDialects []string) bool {
	flag := slices.ContainsFunc(compatibleDialects, func(e string) bool {
		return strings.EqualFold(e, dialect)
	})
	return flag
}

func wrapDriverHooks(drv driver.Driver, drvHooks ...sqlhooks.Hooks) driver.Driver {
	if len(drvHooks) > 0 {
		drv = sqlhooks.Wrap(drv, sqlhooks.Compose(drvHooks...))
	}
	return drv
}
