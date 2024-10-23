package hypersql

import (
	"context"
	"crypto/tls"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	return &dsnConnector{dsn: dsn, drv: drv}, nil
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

	cc := &msdsn.Config{
		Database: name,
		Host:     host,
		Port:     uint64(port),
		User:     user,
		Password: password,
	}

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

	cc.TLSConfig = tlsConfig
	if dialTimeout > 0 {
		cc.DialTimeout = dialTimeout
	}

	if ext, ok := c.Extra.(*SQLServerExtra); ok && ext != nil {

	}

	if err := handleSQLServerParams(params, cc); err != nil {
		return nil, err
	}

	return cc, nil
}

func RawSQLServerDriver() driver.Driver {
	return &mssql.Driver{}
}

func handleSQLServerParams(params ConfigParams, c *msdsn.Config) error {
	var err error
	params.IfNotEmpty(mssqlparams.ConnParams.Database, func(v string) {
		c.Database = v
	})

	if err = params.IfNotEmptyWithErr(mssqlparams.ConnParams.FailOverPort, func(v string) error {
		var err error
		c.FailOverPort, err = strconv.ParseUint(v, 0, 16)
		if err != nil {
			f := "invalid failover port '%v': %v"
			return fmt.Errorf(f, v, err.Error())
		}
		return nil
	}); err != nil {
		return err
	}

	params.IfExists(mssqlparams.ConnParams.HostNameInCertificate, func(v string) {
		c.HostInCertificateProvided = true
	})

	if err = params.IfNotEmptyWithErr(mssqlparams.ConnParams.DisableRetry, func(v string) error {
		var err error
		c.DisableRetry, err = strconv.ParseBool(v)
		if err != nil {
			f := "invalid disableRetry '%s': %s"
			return fmt.Errorf(f, v, err.Error())
		}
		return nil
	}); err != nil {
		return err
	}

	if err = params.IfNotEmptyWithErr(mssqlparams.ConnParams.MultiSubnetFailover, func(v string) error {
		multiSubnetFailover, err := strconv.ParseBool(v)
		if err != nil {
			if strings.EqualFold(v, "Enabled") {
				multiSubnetFailover = true
			} else if strings.EqualFold(v, "Disabled") {
				multiSubnetFailover = false
			} else {
				return fmt.Errorf("invalid multiSubnetFailover value '%v': %v", multiSubnetFailover, err.Error())
			}
		}
		c.MultiSubnetFailover = multiSubnetFailover
		return nil
	}); err != nil {
		return err
	}

	if err = params.IfNotEmptyWithErr(mssqlparams.ConnParams.Encrypt, func(v string) error {
		var encryption msdsn.Encryption = msdsn.EncryptionOff
		v = strings.ToLower(v)
		switch v {
		case "mandatory", "yes", "1", "t", "true":
			encryption = msdsn.EncryptionRequired
		case "disable":
			encryption = msdsn.EncryptionDisabled
		case "strict":
			encryption = msdsn.EncryptionStrict
		case "optional", "no", "0", "f", "false":
			encryption = msdsn.EncryptionOff
		default:
			f := "invalid encrypt '%s'"
			return fmt.Errorf(f, v)
		}
		c.Encryption = encryption
		return nil
	}); err != nil {
		return err
	}

	params.IfNotEmpty(mssqlparams.ConnParams.FailoverPartner, func(v string) {
		c.FailOverPartner = v
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
	params.IfNotEmpty(mssqlparams.ConnParams.DialTimeout, func(v string) {
		tt := cast.ToDuration(v)
		if tt > 0 {
			c.DialTimeout = tt
		}
	})
	params.IfNotEmpty(mssqlparams.ConnParams.ConnectionTimeout, func(v string) {
		tt := cast.ToDuration(v)
		if tt > 0 {
			c.ConnTimeout = tt
		}
	})
	params.IfNotEmpty(mssqlparams.ConnParams.PacketSize, func(v string) {
		psize := cast.ToUint16(v)
		if psize < 512 {
			c.PacketSize = 512
		} else if psize > 32767 {
			c.PacketSize = 32767
		}
	})
	params.IfNotEmpty(mssqlparams.ConnParams.ServerSpn, func(v string) {
		c.ServerSPN = v
	})
	params.IfNotEmpty(mssqlparams.ConnParams.WorkstationID, func(v string) {
		c.Workstation = v
	})
	if err = params.IfNotEmptyWithErr(mssqlparams.ConnParams.ColumnEncryption, func(v string) error {
		columnEncryption, err := strconv.ParseBool(v)
		if err != nil {
			if strings.EqualFold(v, "Enabled") {
				columnEncryption = true
			} else if strings.EqualFold(v, "Disabled") {
				columnEncryption = false
			} else {
				return fmt.Errorf("invalid columnencryption '%v' : %v", columnEncryption, err.Error())
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func IsCompatibleSQLServerDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatibleSQLServerDialects)
}
