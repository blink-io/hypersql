package hypersql

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/qustavo/sqlhooks/v2"
)

var (
	ErrNilConfig = errors.New("[hypersql] config is nil")
)

type DriverHooks []sqlhooks.Hooks

type Config struct {
	Network       string
	Dialect       string
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	Params        ConfigParams
	TLSConfig     *tls.Config
	DialTimeout   time.Duration
	ConnInitSQL   string
	ValidationSQL string
	Loc           *time.Location
	DriverHooks   DriverHooks
	Logger        Logger

	// Connection parameters
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int

	// Enable debug
	Debug bool

	// Database extra config
	Extra any

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
	vv, ok := c.Extra.(Validator)
	if ok {
		return vv.Validate(ctx)
	} else {
		return nil
	}
}

func (c *Config) DBInfo() DBInfo {
	return NewDBInfo(c)
}

func (c *Config) DSN() string {
	return c.dsn
}
