package hypersql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/blink-io/hypersql/mysql/logger"
	mysqlparams "github.com/blink-io/hypersql/mysql/params"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
	"github.com/xo/dburl"
)

var compatibleMySQLDialects = []string{
	DialectMySQL,
	"mysql5",
	"mysql8",
}

func init() {
	dialect := DialectMySQL
	connectors[dialect] = GetMySQLConnector
	dialecters[dialect] = IsCompatibleMySQLDialect
	dsners[dialect] = ToMySQLDSN
}

type MySQLExtra struct {
	ServerPubKey string

	ConnectionAttributes string

	CheckConnLiveness bool

	MultiStatements bool

	InterpolateParams bool
}

var _ Validator = (*MySQLExtra)(nil)

func (c *MySQLExtra) Validate(ctx context.Context) error {
	if c == nil {

	}
	return nil
}

func GetMySQLDSN(dialect string) (Dsner, error) {
	if !IsCompatibleMySQLDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return ToMySQLDSN, nil
}

func GetMySQLDriver(dialect string) (driver.Driver, error) {
	if IsCompatibleMySQLDialect(dialect) {
		return RawMySQLDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetMySQLConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, err := ToMySQLConfig(c)
	if err != nil {
		return nil, err
	}
	dsn := cc.FormatDSN()
	drv := WrapDriver(RawMySQLDriver(), c.DriverWrappers, c.DriverHooks)
	return &dsnConnector{dsn: dsn, drv: drv}, nil
}

func (c *Config) ToMySQL() {
	c.Dialect = DialectMySQL
	c.Port = 3306
}

func ToMySQLConfigFromDSN(dsn string) (*mysql.Config, error) {
	cc, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

func ToMySQLConfigFromURL(url string) (*mysql.Config, error) {
	uu, err := dburl.Parse(url)
	if err != nil {
		return nil, err
	}
	dsn := uu.DSN
	return ToMySQLConfigFromDSN(dsn)
}

func ToMySQLDSN(ctx context.Context, c *Config) (string, error) {
	cc, err := ToMySQLConfig(c)
	if err != nil {
		return "", err
	}
	return cc.FormatDSN(), nil
}

func ToMySQLConfig(c *Config) (*mysql.Config, error) {
	network := c.Transport
	name := c.Name
	host := c.Host
	port := c.Port
	user := c.User
	password := c.Password
	dialTimeout := c.DialTimeout
	tlsConfig := c.TLSConfig
	loc := c.Loc
	params := c.Params
	if loc == nil {
		loc = time.Local
	}
	if params == nil {
		params = make(ConfigParams)
	}

	// Restful TLS Params
	cc := mysql.NewConfig()
	// Put the local timezone because the default value is UTC
	cc.Loc = loc
	// Force to parse to time.Time
	cc.ParseTime = true
	cc.Net = network
	cc.DBName = name
	cc.User = user
	cc.Passwd = password
	if dialTimeout > 0 {
		cc.Timeout = dialTimeout
	}

	if network == "tcp" {
		cc.Addr = net.JoinHostPort(host, cast.ToString(port))
	} else {
		// Otherwise, addr is Unix domain sockets
		cc.Addr = host
	}
	if tlsConfig != nil {
		keyName := mysqlTLSKeyName(name)
		_ = mysql.RegisterTLSConfig(keyName, tlsConfig)
		cc.TLSConfig = keyName
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
		keyName := mysqlTLSKeyName(name)
		_ = mysql.RegisterTLSConfig(keyName, tlscnf)
		cc.TLSConfig = keyName
	}
	if l := c.Logger; l != nil {
		cc.Logger = logger.Logf(func(v ...any) {
			l(fmt.Sprint(v...))
		})
	}

	if ext, ok := c.Extra.(*MySQLExtra); ok && ext != nil {
		cc.ServerPubKey = ext.ServerPubKey
		cc.ConnectionAttributes = ext.ConnectionAttributes
		cc.CheckConnLiveness = ext.CheckConnLiveness
		cc.MultiStatements = ext.MultiStatements
		cc.InterpolateParams = ext.InterpolateParams
	}

	if err := handleMySQLParams(params, cc); err != nil {
		return nil, err
	}

	return cc, nil
}

func IsCompatibleMySQLDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatibleMySQLDialects)
}

func RawMySQLDriver() driver.Driver {
	return &mysql.MySQLDriver{}
}

func mysqlTLSKeyName(name string) string {
	return DialectMySQL + "_" + name
}

func handleMySQLParams(params ConfigParams, cc *mysql.Config) error {
	params.IfNotEmpty(mysqlparams.ConnParams.Collation, func(value string) {
		cc.Collation = value
	})
	return nil
}
