package params

import (
	"strings"
)

// ConnParams specifies connection params.
// Source: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
var ConnParams = connParams{
	Host:                    "host",
	Port:                    "port",
	Hostaddr:                "hostaddr",
	DBName:                  "dbname",
	User:                    "user",
	Passfile:                "passfile",
	ChannelBinding:          "channel_binding",
	RequireAuth:             "require_auth",
	Keepalives:              "keepalives",
	KeepalivesIdle:          "keepalives_idle",
	KeepalivesInterval:      "keepalives_interval",
	KeepalivesCount:         "keepalives_count",
	Replication:             "replication",
	Options:                 "options",
	TCPUserTimeout:          "tcp_user_timeout",
	ApplicationName:         "application_name",
	FallbackApplicationName: "fallback_application_name",
	ClientEncoding:          "client_encoding",
	SSLMode:                 "sslmode",
	SSLKey:                  "sslkey",
	SSLCert:                 "sslcert",
	SSLRootCert:             "sslrootcert",
	SSLPassword:             "sslpassword",
	SSLCertMode:             "sslcertmode",
	SSLCRL:                  "sslcrl",
	SSLCRLDir:               "sslcrldir",
	Service:                 "service",
	ConnectTimeout:          "connect_timeout",
	LoadBalanceHosts:        "load_balance_hosts",
}

// RuntimeParams specifies runtime params.
// Source: https://www.postgresql.org/docs/17/runtime-config-client.html
var RuntimeParams = runtimeParams{
	TimeZone:       "TimeZone",
	ClientEncoding: "client_encoding",
}

type connParams struct {
	// Host specifies name of host to connect to.
	Host string

	// Port specifies port number to connect to at the server host, or socket file name extension for Unix-domain connections.
	Port string

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
	TCPUserTimeout string

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
	ConnectTimeout string

	// Controls the order in which the client tries to connect to the available hosts and addresses.
	LoadBalanceHosts string
}

func (p connParams) Exists(key string) bool {
	switch key := strings.ToLower(key); key {
	case p.Host,
		p.Port,
		p.Hostaddr,
		p.DBName,
		p.User,
		p.Password,
		p.Passfile,
		p.RequireAuth,
		p.Options,
		p.Keepalives,
		p.KeepalivesIdle,
		p.KeepalivesInterval,
		p.KeepalivesCount,
		p.ChannelBinding,
		p.Replication,
		p.TCPUserTimeout,
		p.ApplicationName,
		p.FallbackApplicationName,
		p.ClientEncoding,
		p.SSLMode,
		p.SSLKey,
		p.SSLCert,
		p.SSLRootCert,
		p.SSLPassword,
		p.SSLCertMode,
		p.SSLCRL,
		p.SSLCRLDir,
		p.Service,
		p.ConnectTimeout,
		p.LoadBalanceHosts:
		return true
	default:
		return false
	}
}

type runtimeParams struct {
	TimeZone       string
	ClientEncoding string
}

func (p runtimeParams) Exists(key string) bool {
	switch key := strings.ToLower(key); key {
	case p.TimeZone,
		p.ClientEncoding:
		return true
	default:
		return false
	}
}
