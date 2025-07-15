package postgres

import (
	"github.com/blink-io/hypersql"
	pgparams "github.com/blink-io/hypersql/postgres/params"
	"github.com/spf13/cast"
)

type Config struct {
	// Host specifies name of host to connect to.
	Host string

	// Port specifies port number to connect to at the server host, or socket file name extension for Unix-domain connections.
	Port uint16

	// Numeric IP address of host to connect to.
	Hostaddr string

	// DBName specifies database name.
	DBName string

	// User specifies a PostgreSQL user to connect as.
	User string

	// Password specifies password to be used if the server demands password authentication.
	Password string

	// Passfile specifies the name of the file used to store passwords (see Section 34.16).
	// Defaults to ~/.pgpass or %APPDATA%\postgresql\pgpass.conf on Microsoft Windows.
	// passfile format: hostname:port:database:username:password
	// Source: https://www.postgresql.org/docs/current/libpq-pgpass.html
	Passfile string

	// Specifies the authentication method that the client requires from the server.
	RequireAuth string

	Options string

	// Controls whether client-side TCP keepalives are used.
	Keepalives string

	// Controls the number of seconds of inactivity after which TCP should send a keepalive message to the server.
	KeepalivesIdle string

	// Controls the number of seconds after which a TCP keepalive message that is not acknowledged by the server should be retransmitted.
	KeepalivesInterval string

	// Controls the number of TCP keepalives that can be lost before the client's connection to the server is considered dead.
	KeepalivesCount string

	// This option determines whether the connection should use the replication protocol instead of the normal protocol.
	Replication string

	// Controls the number of milliseconds that transmitted data may remain unacknowledged before a connection is forcibly closed.
	TCPUserTimeout int

	// This option controls the client's use of channel binding.
	ChannelBinding string

	// ApplicationName specifies a value for the application_name configuration parameter.
	ApplicationName string

	// FallbackApplicationName specifies a fallback value for the application_name configuration parameter.
	FallbackApplicationName string

	// ClientEncoding sets the client_encoding configuration parameter for this connection.
	ClientEncoding string

	// SSLMode values: disable|allow|prefer (default)|require|verify-ca|verify-full
	SSLMode string

	// SSLKey specifies the location for the secret key used for the client certificate.
	SSLKey string

	// SSLCert specifies the file name of the client SSL certificate,
	// replacing the default ~/.postgresql/postgresql.crt.
	// This parameter is ignored if an SSL connection is not made.
	SSLCert string

	// SSLRootCert specifies the name of a file containing SSL certificate authority (CA) certificate(s).
	SSLRootCert string

	// SSLPassword specifies the password for the secret key specified in sslkey,
	SSLPassword string

	// SSLCertMode determines whether a client certificate may be sent to the server,
	// and whether the server is required to request one. There are three modes:
	// disable|allow (default)|require
	SSLCertMode string

	// SSLCRL specifies the file name of the SSL server certificate revocation list (CRL).
	SSLCRL string

	// SSLCRLDir specifies the directory name of the SSL server certificate revocation list (CRL).
	SSLCRLDir string

	// Service specifies service name to use for additional parameters.
	Service string

	// ConnectTimeout specifies maximum time to wait while connecting, in seconds (write as a decimal integer, e.g., 10).
	// Zero, negative, or not specified means wait indefinitely.
	ConnectTimeout int

	// Controls the order in which the client tries to connect to the available hosts and addresses.
	LoadBalanceHosts string
}

func (c *Config) ToConfigParams() hypersql.ConfigParams {
	params := hypersql.ConfigParams{
		pgparams.ConnParams.Host:                    c.Host,
		pgparams.ConnParams.Port:                    cast.ToString(c.Port),
		pgparams.ConnParams.Hostaddr:                c.Hostaddr,
		pgparams.ConnParams.DBName:                  c.DBName,
		pgparams.ConnParams.User:                    c.User,
		pgparams.ConnParams.Password:                c.Password,
		pgparams.ConnParams.Passfile:                c.Passfile,
		pgparams.ConnParams.SSLMode:                 c.SSLMode,
		pgparams.ConnParams.SSLKey:                  c.SSLKey,
		pgparams.ConnParams.SSLCert:                 c.SSLCert,
		pgparams.ConnParams.SSLRootCert:             c.SSLRootCert,
		pgparams.ConnParams.SSLPassword:             c.SSLPassword,
		pgparams.ConnParams.SSLCertMode:             c.SSLCertMode,
		pgparams.ConnParams.SSLCRL:                  c.SSLCRL,
		pgparams.ConnParams.SSLCRLDir:               c.SSLCRLDir,
		pgparams.ConnParams.Service:                 c.Service,
		pgparams.ConnParams.ApplicationName:         c.ApplicationName,
		pgparams.ConnParams.FallbackApplicationName: c.FallbackApplicationName,
		pgparams.ConnParams.Replication:             c.Replication,
		pgparams.ConnParams.TCPUserTimeout:          cast.ToString(c.TCPUserTimeout),
		pgparams.ConnParams.ChannelBinding:          c.ChannelBinding,
		pgparams.ConnParams.ClientEncoding:          c.ClientEncoding,
		pgparams.ConnParams.ConnectTimeout:          cast.ToString(c.ConnectTimeout),
		pgparams.ConnParams.LoadBalanceHosts:        c.LoadBalanceHosts,
	}
	return params
}
