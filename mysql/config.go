package mysql

import (
	"github.com/blink-io/hypersql"
	mysqlparams "github.com/blink-io/hypersql/mysql/params"
	"github.com/spf13/cast"
)

type Config struct {
	// Host specifies the host on which the server instance is running.
	Host string

	// Port specifies the TCP/IP network port on which the target MySQL server is listening for connections.
	Port uint16

	// Socket specifies the path to a Unix socket file or the name of a Windows named pipe.
	Socket string

	// Schema specifies the default database for the connection.
	Schema string

	// User specifies the MySQL user account to provide for the authentication process.
	User string

	// Password specifies the password to use for the authentication process.
	Password string

	Loc string

	// SSLMode desired security state for the connection. The following modes are permissible:
	// DISABLED | PREFERRED | REQUIRED | VERIFY_CA | VERIFY_IDENTITY
	SSLMode string

	// SSLCA specifies the path to the X.509 certificate authority file in PEM format.
	SSLCA string

	// SSLCAPath specifies the path to the directory that contains the X.509 certificates authority files in PEM format.
	SSLCAPath string

	// SSLCert specifies the path to the X.509 certificate file in PEM format.
	SSLCert string

	// SSLCRL specifies the path to the file that contains certificate revocation lists in PEM format.
	SSLCRL string

	// SSLCrlpath specifies the path to the directory that contains certificate revocation-list files in PEM format.
	SSLCrlpath string

	// SSLKey specifies the path to the X.509 key file in PEM format.
	SSLKey string

	// TLSVersion specifies the TLS protocols permitted for classic MySQL protocol encrypted connections.
	TLSVersion string

	// AutoMethod specifies the authentication method to use for the connection.
	// The default is AUTO
	AutoMethod string

	Collation string

	Compress bool

	ParseTime bool
}

func (c *Config) ToConfigParams() hypersql.ConfigParams {
	params := hypersql.ConfigParams{
		mysqlparams.ConnParams.Host:       c.Host,
		mysqlparams.ConnParams.Port:       cast.ToString(c.Port),
		mysqlparams.ConnParams.User:       c.User,
		mysqlparams.ConnParams.Password:   c.Password,
		mysqlparams.ConnParams.Schema:     c.Schema,
		mysqlparams.ConnParams.Loc:        c.Loc,
		mysqlparams.ConnParams.SSLMode:    c.SSLMode,
		mysqlparams.ConnParams.SSLCRL:     c.SSLCRL,
		mysqlparams.ConnParams.SSLCert:    c.SSLCert,
		mysqlparams.ConnParams.SSLCA:      c.SSLCA,
		mysqlparams.ConnParams.SSLCAPath:  c.SSLCAPath,
		mysqlparams.ConnParams.SSLKey:     c.SSLKey,
		mysqlparams.ConnParams.SSLCrlpath: c.SSLCrlpath,
		mysqlparams.ConnParams.TLSVersion: c.TLSVersion,
		mysqlparams.ConnParams.AutoMethod: c.AutoMethod,
		mysqlparams.ConnParams.Collation:  c.Collation,
		mysqlparams.ConnParams.Compress:   cast.ToString(c.Compress),
		mysqlparams.ConnParams.ParseTime:  cast.ToString(c.ParseTime),
	}
	return params
}
