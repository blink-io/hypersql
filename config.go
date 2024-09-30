package hypersql

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/qustavo/sqlhooks/v2"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrNilConfig = errors.New("[hypersql] config is nil")
)

type ConfigParams map[string]string

type DriverHooks []sqlhooks.Hooks

type Config struct {
	Network         string
	Dialect         string
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	TLSConfig       *tls.Config
	Params          ConfigParams
	DialTimeout     time.Duration
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnInitSQL     string
	ValidationSQL   string
	DriverHooks     DriverHooks
	Loc             *time.Location
	Debug           bool
	Collation       string
	ClientName      string
	Logger          Logger

	// OpenTelemetry
	WithOTel  bool
	OTelAttrs []attribute.KeyValue

	// Database specific config
	MySQL    *MySQLConfig
	Postgres *PostgresConfig
	SQLite   *SQLiteConfig

	// dsn for internal use
	dsn string
}

func SetupConfig(c *Config) *Config {
	if c == nil {
		c = new(Config)
	}
	c.SetDefaults()
	return c
}

func (c *Config) SetDefaults() {
	if c == nil {
		return
	}

	if len(c.Network) == 0 {
		c.Network = "tcp"
	}
	if c.Loc == nil {
		c.Loc = time.Local
	}
	if len(c.Dialect) > 0 {
		c.Dialect = GetFormalDialect(c.Dialect)
	}
}

func (c *Config) Validate(ctx context.Context) error {
	if c == nil {
		return ErrNilConfig
	}
	d, ok := IsCompatibleDialect(c.Dialect)
	if !ok {
		return ErrUnsupportedDialect
	}
	switch {
	case DialectPostgres == d && c.Postgres != nil:
		return c.Postgres.Validate(ctx)
	case DialectMySQL == d && c.MySQL != nil:
		return c.MySQL.Validate(ctx)
	case DialectSQLite == d && c.SQLite != nil:
		return c.SQLite.Validate(ctx)
	default:
		return ErrUnsupportedDialect
	}
}

func (c *Config) DBInfo() DBInfo {
	return NewDBInfo(c)
}

func (c *Config) DSN() string {
	return c.dsn
}
