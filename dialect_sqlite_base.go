package hypersql

import (
	"github.com/blink-io/hypersql/sqlite"
)

var compatibleSQLiteDialects = []string{
	DialectSQLite,
	"sqlite3",
}

func init() {
	dialect := DialectSQLite
	//drivers[dn] = GetSQLiteDriver
	//dsners[dn] = GetSQLiteDSN
	connectors[dialect] = GetSQLiteConnector

	dialecters[dialect] = IsCompatibleSQLiteDialect
}

func ToSQLiteConfig(c *Config) *sqlite.Config {
	name := c.Name

	cc := &sqlite.Config{
		Name: name,
	}
	return cc
}
