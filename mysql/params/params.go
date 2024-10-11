package params

// Source: https://dev.mysql.com/doc/refman/8.2/en/connecting-using-uri-or-key-value-pairs.html

var ConnParams = struct {
	// Host specifies the host on which the server instance is running.
	Host string

	// Port specifies the TCP/IP network port on which the target MySQL server is listening for connections.
	Port string

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

	Compress string

	ParseTime string
}{
	Host:       "host",
	Port:       "port",
	Socket:     "socket",
	Schema:     "schema",
	User:       "user",
	Password:   "password",
	Loc:        "loc",
	SSLMode:    "ssl-mode",
	SSLCA:      "ssl-ca",
	SSLCAPath:  "ssl-ca-path",
	SSLCert:    "ssl-cert",
	SSLCRL:     "ssl-crl",
	SSLCrlpath: "ssl-crlpath",
	SSLKey:     "ssl-key",
	TLSVersion: "tls-version",
	AutoMethod: "auth-method",
	Collation:  "collation",
	Compress:   "compress",
	ParseTime:  "parseTime",
}
