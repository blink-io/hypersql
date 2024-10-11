package hypersql

import (
	"context"
	"database/sql/driver"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
)

func init() {
	dialect := DialectPostgres
	//drivers[dn] = GetPostgresDriver
	//dsners[dn] = GetPostgresDSN
	connectors[dialect] = GetPostgresConnector

	dialecters[dialect] = IsCompatiblePostgresDialect
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
		return getRawPostgresDriver(), nil
	}
	return nil, ErrUnsupportedDriver
}

func GetPostgresDSN(dialect string) (Dsner, error) {
	if !IsCompatiblePostgresDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return func(ctx context.Context, c *Config) (string, error) {
		cc, err := ToPostgresConfig(c)
		if err != nil {
			return "", err
		}
		dsn := stdlib.RegisterConnConfig(cc)
		return dsn, nil
	}, nil
}

func GetPostgresConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, err := ToPostgresConfig(c)
	if err != nil {
		return nil, err
	}
	c.dsn = stdlib.RegisterConnConfig(cc)
	drv := wrapDriver(getRawPostgresDriver(), c.DriverWrappers, c.DriverHooks)
	return &dsnConnector{dsn: c.dsn, driver: drv}, nil
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

func ToPostgresConfig(c *Config) (*pgx.ConnConfig, error) {
	name := c.Name
	host := c.Host
	port := c.Port
	user := c.User
	password := c.Password
	dialTimeout := c.DialTimeout
	tlsConfig := c.TLSConfig
	params := c.Params
	if params == nil {
		params = make(map[string]string)
	}

	pgcc, err := pgconn.ParseConfig("")
	if err != nil {
		return nil, err
	}

	pgcc.Database = name
	pgcc.Host = host
	pgcc.Port = uint16(port)
	pgcc.User = user
	pgcc.Password = password
	pgcc.TLSConfig = tlsConfig
	pgcc.RuntimeParams = handlePostgresParams(params)
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

func getRawPostgresDriver() driver.Driver {
	// Notes: Unable to invoke &stdlib.Driver{} directly.
	// Because the "configs" field inside the driver is not initialized.
	return stdlib.GetDefaultDriver()
}

func handlePostgresParams(params map[string]string) map[string]string {
	newParams := make(map[string]string)
	for k, v := range params {
		newParams[k] = v
	}
	return newParams
}

func IsCompatiblePostgresDialect(dialect string) bool {
	return isCompatibleDialectIn(dialect, compatiblePostgresDialects)
}
