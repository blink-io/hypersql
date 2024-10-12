package hypersql

import (
	"context"
	"crypto/tls"
	"errors"
	"time"
)

const (
	XParamDialTimeout   = "x-dial-timeout"
	XParamConnInitSQL   = "x-conn-init-sql"
	XParamValidationSQL = "x-conn-validation-sql"

	XParamConnMaxLifetime = "x-conn-max-lifetime"
	XParamConnMaxIdleTime = "x-conn-max-idle-time"
	XParamMaxOpenConns    = "x-max-open-conns"
	XParamMaxIdleConns    = "x-max-idle-conns"
)

var (
	ErrNilConfig = errors.New("[hypersql] config is nil")
)

type Config struct {
	Transport     string         `json:"transport" yaml:"transport" toml:"transport"`
	Dialect       string         `json:"dialect" yaml:"dialect"  toml:"dialect"`
	Host          string         `json:"host" yaml:"host" toml:"host"`
	Port          int            `json:"port" yaml:"port" toml:"port"`
	Name          string         `json:"name" yaml:"name" toml:"name"`
	User          string         `json:"user" yaml:"user" toml:"user"`
	Password      string         `json:"password" yaml:"password" toml:"password"`
	Params        ConfigParams   `json:"params" yaml:"params" toml:"params"`
	DialTimeout   time.Duration  `json:"dial_timeout" yaml:"dial_timeout" toml:"dial_timeout"`
	ConnInitSQL   string         `json:"conn_init_sql" yaml:"conn_init_sql" toml:"conn_init_sql"`
	ValidationSQL string         `json:"validation_sql" yaml:"validation_sql" toml:"validation_sql"`
	Loc           *time.Location `json:"loc" yaml:"loc" toml:"loc"`
	Logger        Logger         `json:"-" yaml:"-" toml:"-"`

	TLSCert   *TLSCert    `json:"tls_cert" yaml:"tls_cert" toml:"tls_cert"`
	TLSConfig *tls.Config `json:"-" yaml:"-" toml:"-"`

	// Driver related
	DriverHooks    DriverHooks
	DriverWrappers DriverWrappers

	AfterHandlers AfterHandlers

	// Connection parameters
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime" toml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time" toml:"conn_max_idle_time"`
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns" toml:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns" toml:"max_idle_conns"`

	// Enable debug
	Debug bool

	// Database extra config
	Extra any
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

	if len(c.Transport) == 0 {
		c.Transport = "tcp"
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
