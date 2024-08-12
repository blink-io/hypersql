package hypersql

import (
	"context"
	"database/sql/driver"

	pgparams "github.com/blink-io/hypersql/postgres/params"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func init() {
	d := DialectPostgres
	//drivers[dn] = GetPostgresDriver
	//dsners[dn] = GetPostgresDSN
	connectors[d] = GetPostgresConnector
}

var compatiblePostgresDialects = []string{
	DialectPostgres,
	"postgresql",
	"pg",
	"pgx",
}

type PostgresConfig struct {
	UsePool bool

	DialFunc pgconn.DialFunc

	AfterConnect pgconn.AfterConnectFunc

	ValidateConnect pgconn.ValidateConnectFunc

	OnNotice pgconn.NoticeHandler

	OnNotification pgconn.NotificationHandler

	Tracer pgx.QueryTracer

	StatementCacheCapacity int

	DescriptionCacheCapacity int
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
		cc, err := ToPGXConfig(c)
		if err != nil {
			return "", err
		}
		dsn := stdlib.RegisterConnConfig(cc)
		return dsn, nil
	}, nil
}

func GetPostgresConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, err := ToPGXConfig(c)
	if err != nil {
		return nil, err
	}

	var usePool = false
	if c.Postgres != nil {
		usePool = c.Postgres.UsePool
	}
	if usePool {
		ppc, errv := pgxpool.ParseConfig("")
		if errv != nil {
			return nil, errv
		}
		ppc.ConnConfig = cc
		pool, errp := pgxpool.NewWithConfig(ctx, ppc)
		if errp != nil {
			return nil, errp
		}
		c.dsn = ppc.ConnString()
		return stdlib.GetPoolConnector(pool), nil
	} else {
		c.dsn = stdlib.RegisterConnConfig(cc)
		drv := wrapDriverHooks(getRawPostgresDriver(), c.DriverHooks...)
		return &dsnConnector{dsn: c.dsn, driver: drv}, nil
	}
}

func ToPGXConfig(c *Config) (*pgx.ConnConfig, error) {
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
	if len(c.ClientName) > 0 {
		params[pgparams.ApplicationName] = c.ClientName
	}
	if len(c.Collation) > 0 {
		params[pgparams.ClientEncoding] = c.Collation
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

	if pc := c.Postgres; pc != nil {
		if pc.Tracer != nil {
			cc.Tracer = pc.Tracer
		}
		if pc.StatementCacheCapacity > 0 {
			cc.StatementCacheCapacity = pc.StatementCacheCapacity
		}
		if pc.DescriptionCacheCapacity > 0 {
			cc.DescriptionCacheCapacity = pc.DescriptionCacheCapacity
		}
	}

	return cc, nil
}

func ToPGXPoolConfig(c *Config) (*pgxpool.Config, error) {
	ppc, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	cfg, err := ToPGXConfig(c)
	if err != nil {
		return nil, err
	}

	ppc.ConnConfig = cfg
	return ppc, nil
}

func getRawPostgresDriver() driver.Driver {
	// Notes: Unable to invoke &stdlib.Driver{} directly.
	// Because the "configs" field inside the driver is not initialized.
	return stdlib.GetDefaultDriver()
}

func ValidatePostgresConfig(c *Config) error {
	return nil
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
