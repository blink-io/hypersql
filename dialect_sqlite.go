package hypersql

import (
	"context"
)

var compatibleSQLiteDialects = []string{
	DialectSQLite,
	"sqlite3",
}

func init() {
	d := DialectSQLite
	//drivers[dn] = GetSQLiteDriver
	//dsners[dn] = GetSQLiteDSN
	connectors[d] = GetSQLiteConnector
}

type SQLiteConfig struct {
}

func (c *SQLiteConfig) Validate(ctx context.Context) error {
	if c == nil {
		return ErrNilConfig
	}
	return nil
}

func toSQLiteDSN(c *Config) string {
	dsn := c.Host
	return dsn
}
