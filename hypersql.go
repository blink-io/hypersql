package hypersql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"runtime"

	"github.com/spf13/cast"
)

const (
	// DialectMySQL defines MySQL dialect
	DialectMySQL = "mysql"
	// DialectPostgres defines PostgreSQL dialect
	DialectPostgres = "postgres"
	// DialectSQLite defines SQLite dialect. We only support SQLite3 above.
	DialectSQLite = "sqlite"
)

var (
	ErrUnsupportedDialect = errors.New("unsupported dialect")

	ErrUnsupportedDriver = errors.New("unsupported driver")

	InvalidConfig = errors.New("invalid config")
)

type (
	Pinger interface {
		PingContext(ctx context.Context) error
	}

	IDB interface {
		Begin() (*sql.Tx, error)
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	}

	IDBExt interface {
		WithSqlDB

		WithDBInfo

		HealthChecker
	}

	WithSqlDB interface {
		SqlDB() *sql.DB
	}

	WithDBInfo interface {
		DBInfo() DBInfo
	}

	HealthChecker interface {
		HealthCheck(context.Context) error
	}

	DBInfo struct {
		Name    string
		Dialect string
	}

	Validator interface {
		Validate(context.Context) error
	}
)

func NewSqlDB(c *Config) (*sql.DB, error) {
	dialect := GetFormalDialect(c.Dialect)
	connFn := GetConnector(dialect)
	if connFn == nil {
		return nil, ErrUnsupportedDialect
	}

	ctx := context.Background()
	conn, err := connFn(ctx, c)
	if err != nil {
		return nil, err
	}

	var db *sql.DB
	if c.WithOTel {
		otelOps := []OTelOption{
			OTelDBHostPort(hostPortToAddr(c.Host, c.Port)),
			OTelDBName(c.Name),
			OTelDBSystem(dialect),
			OTelReportDBStats(),
			OTelAttrs(c.OTelAttrs...),
		}
		db = otelOpenDB(conn, otelOps...)
	} else {
		db = otelWrapper(sql.OpenDB)(conn)
	}

	// Do ping check
	if err := DoPingContext(ctx, db); err != nil {
		return nil, err
	}

	var doExec = func(action string, sql string) error {
		if len(sql) > 0 {
			if _, err := db.Exec(sql); err != nil {
				return fmt.Errorf("unable to exec sql for [%s]: %s, reason: %s", action, sql, err)
			}
		}
		return nil
	}

	if err := doExec("connection initialization", c.ConnInitSQL); err != nil {
		return nil, err
	}

	if err := doExec("validation", c.ValidationSQL); err != nil {
		return nil, err
	}

	// Reference: https://bun.uptrace.dev/guide/running-bun-in-production.html
	maxIdleConns := c.MaxIdleConns
	maxOpenConns := c.MaxOpenConns
	connMaxLifetime := c.ConnMaxLifetime
	connMaxIdleTime := c.ConnMaxIdleTime
	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
	} else {
		// TODO In Docker how we should do?
		maxOpenConns = 4 * runtime.GOMAXPROCS(0)
		db.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns)
	} else {
		db.SetMaxIdleConns(maxOpenConns)
	}
	if connMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(connMaxIdleTime)
	}
	if connMaxLifetime > 0 {
		db.SetConnMaxLifetime(connMaxLifetime)
	}

	return db, nil
}

func NewDBInfo(c *Config) DBInfo {
	return DBInfo{
		Name:    c.Name,
		Dialect: c.Dialect,
	}
}

func hostPortToAddr(host string, port int) string {
	return net.JoinHostPort(host, cast.ToString(port))
}
