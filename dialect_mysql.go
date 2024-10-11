package hypersql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net"
	"time"

	"github.com/blink-io/hypersql/mysql/logger"
	mysqlparams "github.com/blink-io/hypersql/mysql/params"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
)

var compatibleMySQLDialects = []string{
	DialectMySQL,
	"mysql5",
	"mysql8",
}

func init() {
	dialect := DialectMySQL
	//drivers[d] = GetMySQLDriver
	//dsners[d] = GetMySQLDSN
	connectors[dialect] = GetMySQLConnector

	dialecters[dialect] = IsCompatibleMySQLDialect
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
	return func(ctx context.Context, c *Config) (string, error) {
		cc, err := ToMySQLConfig(c)
		if err != nil {
			return "", err
		}
		dsn := cc.FormatDSN()
		return dsn, nil
	}, nil
}

func GetMySQLDriver(dialect string) (driver.Driver, error) {
	if IsCompatibleMySQLDialect(dialect) {
		return getRawMySQLDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetMySQLConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, err := ToMySQLConfig(c)
	if err != nil {
		return nil, err
	}
	c.dsn = cc.FormatDSN()
	drv := wrapDriver(getRawMySQLDriver(), c.DriverWrappers, c.DriverHooks)
	return &dsnConnector{dsn: c.dsn, driver: drv}, nil
}

func ToMySQLConfig(c *Config) (*mysql.Config, error) {
	network := c.Network
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
		params = make(map[string]string)
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
	// TODO Do we need to check them?
	cc.Params = handleMySQLParams(params)
	params.IfNotEmpty(mysqlparams.Collation, func(value string) {
		cc.Collation = value
	})
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

	return cc, nil
}

func IsCompatibleMySQLDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatibleMySQLDialects)
}

func getRawMySQLDriver() driver.Driver {
	return &mysql.MySQLDriver{}
}

func mysqlTLSKeyName(name string) string {
	return DialectMySQL + "_" + name
}

func handleMySQLParams(params map[string]string) map[string]string {
	newParams := make(map[string]string)
	for k, v := range params {
		newParams[k] = v
	}
	return newParams
}
