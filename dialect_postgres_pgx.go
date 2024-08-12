package hypersql

import (
	"context"
	"database/sql/driver"
	"log/slog"

	pgparams "github.com/blink-io/hypersql/postgres/params"
	pgxslog "github.com/blink-io/hypersql/postgres/pgx/logger/slog"
	pgxotel "github.com/blink-io/hypersql/postgres/pgx/tracer/otel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/spf13/cast"
)

func GetPostgresDSN(dialect string) (Dsner, error) {
	if !IsCompatiblePostgresDialect(dialect) {
		return nil, ErrUnsupportedDialect
	}
	return func(ctx context.Context, c *Config) (string, error) {
		cc, _, err := ToPGXConfig(c)
		if err != nil {
			return "", err
		}
		dsn := stdlib.RegisterConnConfig(cc)
		return dsn, nil
	}, nil
}

func GetPostgresConnector(ctx context.Context, c *Config) (driver.Connector, error) {
	cc, aopt, err := ToPGXConfig(c)
	if err != nil {
		return nil, err
	}
	if aopt.usePool {
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

func ToPGXConfig(c *Config) (*pgx.ConnConfig, *PostgresOptions, error) {
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
		return nil, nil, err
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
		return nil, nil, err
	}

	opts := AdditionsToPostgresOptions(c.Additions)

	cc.Config = *pgcc
	traceLogLevel := tracelog.LogLevelInfo
	if c.Debug {
		traceLogLevel = tracelog.LogLevelDebug
	} else {
		if lv, errv := tracelog.LogLevelFromString(opts.loglevel); errv == nil {
			traceLogLevel = lv
		}
	}

	if opts.trace == PostgresTraceOTel {
		cc.Tracer = pgxotel.NewTracer()
	} else {
		var tlogger tracelog.Logger
		if l := c.Logger; l != nil {
			tlogger = tracelog.LoggerFunc(doLoggerFunc(l))
		} else {
			tlogger = pgxslog.NewLogger(slog.Default())
		}
		cc.Tracer = &tracelog.TraceLog{Logger: tlogger, LogLevel: traceLogLevel}
	}
	return cc, opts, nil
}

func ToPGXPoolConfig(c *Config) (*pgxpool.Config, *PostgresOptions, error) {
	ppc, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, nil, err
	}

	cfg, opts, err := ToPGXConfig(c)
	if err != nil {
		return nil, nil, err
	}

	ppc.ConnConfig = cfg
	return ppc, opts, nil
}

func AdditionsToPostgresOptions(adds map[string]string) *PostgresOptions {
	opts := new(PostgresOptions)
	if adds != nil {
		opts.trace = adds[AdditionTrace]
		opts.tracelog = adds[AdditionTracelog]
		opts.usePool = cast.ToBool(adds[AdditionUsePool])
	}
	return opts
}

func getRawPostgresDriver() driver.Driver {
	// Notes: Unable to invoke &stdlib.Driver{} directly.
	// Because the "configs" field inside the driver is not initialized.
	return stdlib.GetDefaultDriver()
}
