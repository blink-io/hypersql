package hypersql

import (
	"context"
	"crypto/tls"
	"database/sql/driver"
	"errors"

	mssqlparams "github.com/blink-io/hypersql/sqlserver/params"
	"github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/msdsn"
	"github.com/spf13/cast"
)

func init() {
	dialect := DialectSQLServer
	connectors[dialect] = GetSQLServerConnector
	dialecters[dialect] = IsCompatibleSQLServerDialect
	dsners[dialect] = ToSQLServerDSN
}

var compatibleSQLServerDialects = []string{
	DialectSQLServer,
	"mssql",
	"mssqldb",
}

type SQLServerExtra struct {
}

var _ Validator = (*SQLServerExtra)(nil)

func (c *SQLServerExtra) Validate(ctx context.Context) error {
	if c == nil {

	}
	return nil
}

func GetSQLServerDriver(dialect string) (driver.Driver, error) {
	if IsCompatibleSQLServerDialect(dialect) {
		return RawSQLServerDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetSQLServerDSN(dialect string) (Dsner, error) {
	if !IsCompatibleSQLServerDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return ToSQLServerDSN, nil
}

func GetSQLServerConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, err := ToSQLServerConfig(c)
	if err != nil {
		return nil, err
	}
	dsn := cc.URL().String()
	drv := WrapDriver(RawSQLServerDriver(), c.DriverWrappers, c.DriverHooks)
	return &dsnConnector{dsn: dsn, dri: drv}, nil
}

func (c *Config) ToSQLServer() {
	c.Dialect = DialectSQLServer
	c.Port = 1434
}

// ToSQLServerConfigFromDSN converts the config to SQLServer DSN config.
func ToSQLServerConfigFromDSN(dsn string) (*msdsn.Config, error) {
	cc, err := msdsn.Parse(dsn)
	if err != nil {
		return nil, err
	}
	return &cc, err
}

// ToSQLServerConfigFromURL converts the config to SQLServer URL config.
func ToSQLServerConfigFromURL(url string) (*msdsn.Config, error) {
	cc, err := msdsn.Parse(url)
	if err != nil {
		return nil, err
	}
	return &cc, err
}

// ToSQLServerDSN converts the config to SQLServer DSN string.
func ToSQLServerDSN(ctx context.Context, c *Config) (string, error) {
	cc, err := ToSQLServerConfig(c)
	if err != nil {
		return "", err
	}
	return cc.URL().String(), nil
}

func ToSQLServerConfig(c *Config) (*msdsn.Config, error) {
	var tlsConfig *tls.Config
	name := c.Name
	host := c.Host
	port := c.Port
	user := c.User
	password := c.Password
	dialTimeout := c.DialTimeout
	params := c.Params
	if params == nil {
		params = make(ConfigParams)
	}

	cc := &msdsn.Config{}

	if c.TLSConfig != nil {
		tlsConfig = c.TLSConfig
	} else if c.TLSCert != nil {
		tlscnf, err := CreateClientTLSConfig(
			c.TLSCert.CAFile,
			c.TLSCert.CAOptional,
			c.TLSCert.CertFile,
			c.TLSCert.KeyFile,
			c.TLSCert.InsecureSkipVerify,
		)
		if err != nil {
			return nil, errors.New("invalid ca file or key file")
		}
		tlsConfig = tlscnf
	}

	cc.Database = name
	cc.Host = host
	cc.Port = uint64(port)
	cc.User = user
	cc.Password = password
	cc.TLSConfig = tlsConfig
	if dialTimeout > 0 {
		cc.DialTimeout = dialTimeout
	}

	if err := processSQLServerParams(params, cc); err != nil {
		return nil, err
	}

	if ext, ok := c.Extra.(*SQLServerExtra); ok && ext != nil {

	}

	return cc, nil
}

func RawSQLServerDriver() driver.Driver {
	return &mssql.Driver{}
}

func processSQLServerParams(params ConfigParams, c *msdsn.Config) error {
	params.IfNotEmpty(mssqlparams.ConnParams.Database, func(v string) {
		c.Database = v
	})
	params.IfNotEmpty(mssqlparams.ConnParams.AppName, func(v string) {
		c.AppName = v
	})
	params.IfNotEmpty(mssqlparams.ConnParams.KeepAlive, func(v string) {
		tt := cast.ToDuration(v)
		if tt > 0 {
			c.KeepAlive = tt
		}
	})
	params.IfNotEmpty(mssqlparams.ConnParams.ServerSpn, func(v string) {
		c.ServerSPN = v
	})
	params.IfNotEmpty(mssqlparams.ConnParams.WorkstationID, func(v string) {
		c.Workstation = v
	})
	return nil
}

func IsCompatibleSQLServerDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatibleSQLServerDialects)
}
