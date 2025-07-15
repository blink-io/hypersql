//go:build sqlite

package hypersql

import (
	"context"
	"database/sql/driver"

	"github.com/blink-io/hypersql/sqlite"
	"github.com/xo/dburl"
)

var compatibleSQLiteDialects = []string{
	DialectSQLite,
	DialectSQLite3,
}

func init() {
	//drivers[dn] = GetSQLiteDriver
	//dsners[dn] = GetSQLiteDSN
	connectors[DialectSQLite] = GetSQLiteConnector
	dialecters[DialectSQLite] = IsCompatibleSQLiteDialect
	dsners[DialectSQLite] = ToSQLiteDSN

	connectors[DialectSQLite3] = GetSQLiteConnector
	dialecters[DialectSQLite3] = IsCompatibleSQLiteDialect
	dsners[DialectSQLite3] = ToSQLiteDSN
}

func GetSQLiteDSN(dialect string) (Dsner, error) {
	if !IsCompatibleSQLiteDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return ToSQLiteDSN, nil
}

func IsCompatibleSQLiteDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatibleSQLiteDialects)
}

func GetSQLiteDriver(dialect string) (driver.Driver, error) {
	if IsCompatibleSQLiteDialect(dialect) {
		return RawSQLiteDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetSQLiteConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, err := ToSQLiteConfig(c)
	if err != nil {
		return nil, err
	}
	dsn := cc.FormatDSN()
	drv := WrapDriver(RawSQLiteDriver(), c.DriverWrappers, c.DriverHooks)
	return &dsnConnector{dsn: dsn, drv: drv}, nil
}

func (c *Config) ToSQLite() {
	c.Dialect = DialectSQLite
}

func ToSQLiteConfigFromDSN(ctx context.Context, dsn string) (*sqlite.Config, error) {
	cc, err := sqlite.ParseDSN(dsn)
	return cc, err
}

func ToSQLiteConfigFromURL(ctx context.Context, url string) (*sqlite.Config, error) {
	uu, err := dburl.Parse(url)
	if err != nil {
		return nil, err
	}
	cc := &sqlite.Config{
		Name: uu.Path,
	}
	params := make(map[string]string)
	for k := range uu.Query() {
		params[k] = uu.Query().Get(k)
	}
	if err := cc.HandleParams(params); err != nil {
		return nil, err
	}
	return cc, nil
}

func ToSQLiteDSN(ctx context.Context, c *Config) (string, error) {
	cc, err := ToSQLiteConfig(c)
	if err != nil {
		return "", err
	}
	return cc.FormatDSN(), nil
}

func ToSQLiteConfig(c *Config) (*sqlite.Config, error) {
	params := c.Params
	cc := &sqlite.Config{
		Name: c.Name,
	}
	if err := cc.HandleParams(params); err != nil {
		return nil, err
	}
	return cc, nil
}
