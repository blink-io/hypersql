package hypersql

import (
	"context"
	"crypto/tls"
	"database/sql/driver"
	"errors"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cast"
)

func init() {
	dialect := DialectPostgres
	connectors[dialect] = GetPostgresConnector
	dialecters[dialect] = IsCompatiblePostgresDialect
	dsners[dialect] = ToPostgresDSN
}

var compatiblePostgresDialects = []string{
	DialectPostgres,
	"postgresql",
	"pg",
	"pgx",
}

type PostgresExtra struct {
	DialFunc pgconn.DialFunc

	AfterConnect pgconn.AfterConnectFunc

	ValidateConnect pgconn.ValidateConnectFunc

	OnNotice pgconn.NoticeHandler

	OnNotification pgconn.NotificationHandler

	Tracers []pgx.QueryTracer

	StatementCacheCapacity int

	DescriptionCacheCapacity int
}

var _ Validator = (*PostgresExtra)(nil)

func (c *PostgresExtra) Validate(ctx context.Context) error {
	if c == nil {

	}
	return nil
}

func GetPostgresDriver(dialect string) (driver.Driver, error) {
	if IsCompatiblePostgresDialect(dialect) {
		return RawPostgresDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetPostgresDSN(dialect string) (Dsner, error) {
	if !IsCompatiblePostgresDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return ToPostgresDSN, nil
}

func GetPostgresConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, err := ToPostgresConfig(c)
	if err != nil {
		return nil, err
	}
	dsn := stdlib.RegisterConnConfig(cc)
	drv := WrapDriver(RawPostgresDriver(), c.DriverWrappers, c.DriverHooks)
	return &dsnConnector{dsn: dsn, drv: drv}, nil
}

func (c *Config) ToPostgres() {
	c.Dialect = DialectPostgres
	c.Port = 5432
}

// ToPostgresConfigFromDSN
// See: https://github.com/jackc/pgx/blob/v5.4.0/conn.go#L101
// See: https://github.com/lib/pq/blob/v1.10.9/url.go#L32
func ToPostgresConfigFromDSN(dsn string) (*pgx.ConnConfig, error) {
	cc, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	return cc, err
}

func ToPostgresConfigFromURL(url string) (*pgx.ConnConfig, error) {
	cc, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	return cc, err
}

// ToPostgresDSN converts the config to PostgreSQL DSN string.
// See https://github.com/lib/pq/blob/v1.10.9/url.go#L32
func ToPostgresDSN(ctx context.Context, c *Config) (string, error) {
	var kvs []string
	escaper := strings.NewReplacer(`'`, `\'`, `\`, `\\`)
	accrue := func(k, v string) {
		if v != "" {
			kvs = append(kvs, k+"='"+escaper.Replace(v)+"'")
		}
	}

	accrue("user", c.User)
	accrue("password", c.Password)
	accrue("host", c.Host)
	accrue("port", cast.ToString(c.Port))
	accrue("dbname", c.Name)

	q := c.Params
	for k, v := range q {
		accrue(k, v)
	}

	sort.Strings(kvs) // Makes testing easier (not a performance concern)
	return strings.Join(kvs, " "), nil
}

func ToPostgresConfig(c *Config) (*pgx.ConnConfig, error) {
	var tlsConfig *tls.Config
	name := c.Name
	host := c.Host
	port := c.Port
	user := c.User
	password := c.Password
	dialTimeout := c.DialTimeout
	params := c.Params
	if params == nil {
		params = make(map[string]string)
	}

	pgcc, err := pgconn.ParseConfig("")
	if err != nil {
		return nil, err
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

	pgcc.Database = name
	pgcc.Host = host
	pgcc.Port = uint16(port)
	pgcc.User = user
	pgcc.Password = password
	pgcc.RuntimeParams = handlePostgresParams(params)
	pgcc.TLSConfig = tlsConfig
	if dialTimeout > 0 {
		pgcc.ConnectTimeout = dialTimeout
	}

	cc, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}
	cc.Config = *pgcc

	if ext, ok := c.Extra.(*PostgresExtra); ok && ext != nil {
		if tracers := ext.Tracers; len(tracers) > 0 {
			cc.Tracer = multitracer.New(tracers...)
		}
		if ext.StatementCacheCapacity > 0 {
			cc.StatementCacheCapacity = ext.StatementCacheCapacity
		}
		if ext.DescriptionCacheCapacity > 0 {
			cc.DescriptionCacheCapacity = ext.DescriptionCacheCapacity
		}
	}

	return cc, nil
}

func RawPostgresDriver() driver.Driver {
	// Notes: Unable to invoke &stdlib.Driver{} directly.
	// Because the "configs" field inside the drv is not initialized.
	return stdlib.GetDefaultDriver()
}

func handlePostgresParams(params ConfigParams) map[string]string {
	newParams := make(map[string]string)
	for k, v := range params {
		newParams[k] = v
	}
	return newParams
}

func IsCompatiblePostgresDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatiblePostgresDialects)
}
