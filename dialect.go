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

	dialecters = make(map[string]func(string) bool)
)

func GetConnector(dialect string) ConnectorFunc {
	return connectors[dialect]
}

func GetFormalDialect(dialect string) string {
	if d, ok := IsCompatibleDialect(dialect); ok {
		return d
	}
	return ""
}

// IsCompatibleDialect checks
func IsCompatibleDialect(dialect string) (string, bool) {
	for k, v := range dialecters {
		if v(dialect) {
			return k, true
		}
	}
	return "", false
}

func isCompatibleDialectIn(dialect string, compatibleDialects []string) bool {
	flag := slices.ContainsFunc(compatibleDialects, func(e string) bool {
		return strings.EqualFold(e, dialect)
	})
	return flag
}

func wrapDriver(drv driver.Driver, wrappers DriverWrappers, hooks DriverHooks) driver.Driver {
	return wrapDriverHooks(wrapDriverWrappers(drv, wrappers...), hooks...)
}

func wrapDriverHooks(drv driver.Driver, drvHooks ...sqlhooks.Hooks) driver.Driver {
	if len(drvHooks) > 0 {
		drv = sqlhooks.Wrap(drv, sqlhooks.Compose(drvHooks...))
	}
	return drv
}

func wrapDriverWrappers(drv driver.Driver, wrappers ...DriverWrapper) driver.Driver {
	for _, w := range wrappers {
		drv = w(drv)
	}
	return drv
}
