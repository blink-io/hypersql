package sqlite

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/blink-io/hypersql"
	sqliteparams "github.com/blink-io/hypersql/sqlite/params"
	"github.com/spf13/cast"
)

const (
	CryptSHA1    string = "SHA1"
	CryptSSHA1   string = "SSHA1"
	CryptSHA256  string = "SHA256"
	CryptSSHA256 string = "SSHA256"
	CryptSHA384  string = "SHA384"
	CryptSSHA384 string = "SSHA384"
	CryptSHA512  string = "SHA512"
	CryptSSHA512 string = "SSHA512"

	AutoVacuumNone        string = "none"
	AutoVacuumFull        string = "full"
	AutoVacuumIncremental string = "incremental"

	BoolTrueOn   = "on"
	BoolTrueYes  = "yes"
	BoolTrue1    = "1"
	BoolTrueTrue = "true"

	BoolFalseOff   = "off"
	BoolFalseNo    = "no"
	BoolFalse0     = "0"
	BoolFalseFalse = "false"

	JournalDelete   string = "DELETE"
	JournalTruncate string = "TRUNCATE"
	JournalPersist  string = "PERSIST"
	JournalMemory   string = "MEMORY"
	JournalWal      string = "WAL"
	JournalOff      string = "OFF"

	LockingModeNormal    string = "NORMAL"
	LockingModeExclusive string = "EXCLUSIVE"

	ModeRO     string = "ro"
	ModeRW     string = "rw"
	ModeRWC    string = "rwc"
	ModeMemory string = "memory"

	MutexFull string = "full"
	MutexNo   string = "no"

	SyncFull   string = "full"
	SyncNormal string = "normal"
	SyncOff    string = "off"
	SyncExtra  string = "extra"

	LocAuto string = "auto"

	TxLockImmediate string = "immediate"
	TxLockExclusive string = "exclusive"
	TxLockDeferred  string = "deferred"

	SecureDeleteFast = "FAST"

	CacheShared  string = "shared"
	CachePrivate string = "private"
)

var (
	ErrSynonymousParam = errors.New("sqlite: synonymous param")
)

// Config is the string for the SQLite database.
// See https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string
// See https://pkg.go.dev/github.com/mattn/go-sqlite3#SQLiteDriver
type Config struct {
	Name                   string
	Auth                   bool
	AuthUser               string
	AuthPass               string
	AuthCrypt              string
	AuthSalt               string
	AutoVacuum             string
	BusyTimeout            int
	Cache                  string
	CacheSize              int
	CaseSensitiveLike      bool
	DeferForeignKeys       bool
	ForeignKeys            bool
	IgnoreCheckConstraints bool
	Immutable              bool
	JournalMode            string
	Loc                    string
	LockingMode            string
	Mode                   string
	Mutex                  string
	QueryOnly              bool
	RecursiveTriggers      bool
	SecureDelete           string
	Sync                   string
	TxLock                 string
	WritableSchema         bool
}

func (c *Config) FormatDSN() string {
	return FormatDSN(c)
}

func FormatDSN(c *Config) string {
	var buf bytes.Buffer
	var isFirst = true
	var accrue = func(key string, value any) {
		if isFirst {
			buf.WriteString("?")
			isFirst = false
		} else {
			buf.WriteString("&")
		}
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(cast.ToString(value))
	}

	buf.WriteString(c.Name)

	if c.Immutable {
		accrue(sqliteparams.ConnParams.Immutable, c.Immutable)
	}
	if c.QueryOnly {
		accrue(sqliteparams.ConnParams.QueryOnly, c.QueryOnly)
	}
	if c.Auth {
		accrue(sqliteparams.ConnParams.Auth, "")
	}
	if len(c.AuthUser) > 0 {
		accrue(sqliteparams.ConnParams.AuthUser, c.AuthUser)
	}
	if len(c.AuthPass) > 0 {
		accrue(sqliteparams.ConnParams.AuthPass, c.AuthPass)
	}
	if len(c.AuthCrypt) > 0 {
		accrue(sqliteparams.ConnParams.AuthCrypt, c.AuthCrypt)
	}
	if len(c.AuthSalt) > 0 {
		accrue(sqliteparams.ConnParams.AuthSalt, c.AuthSalt)
	}
	if len(c.Cache) > 0 {
		accrue(sqliteparams.ConnParams.Cache, c.Cache)
	}
	if len(c.Mode) > 0 {
		accrue(sqliteparams.ConnParams.Mode, c.Mode)
	}
	if len(c.Mutex) > 0 {
		accrue(sqliteparams.ConnParams.Mutex, c.Mutex)
	}
	if c.BusyTimeout > 0 {
		accrue(sqliteparams.ConnParams.BusyTimeout, c.BusyTimeout)
	}
	if len(c.AutoVacuum) > 0 {
		accrue(sqliteparams.ConnParams.AutoVacuum, c.AutoVacuum)
	}
	if len(c.JournalMode) > 0 {
		accrue(sqliteparams.ConnParams.JournalMode, c.JournalMode)
	}
	if len(c.LockingMode) > 0 {
		accrue(sqliteparams.ConnParams.LockingMode, c.LockingMode)
	}
	if len(c.SecureDelete) > 0 {
		accrue(sqliteparams.ConnParams.SecureDelete, c.SecureDelete)
	}
	if len(c.Loc) > 0 {
		accrue(sqliteparams.ConnParams.Loc, c.Loc)
	}
	if len(c.Sync) > 0 {
		accrue(sqliteparams.ConnParams.Sync, c.Sync)
	}
	if len(c.TxLock) > 0 {
		accrue(sqliteparams.ConnParams.TxLock, c.TxLock)
	}
	if c.CacheSize > 0 {
		accrue(sqliteparams.ConnParams.CacheSize, c.CacheSize)
	}

	// Append other params when it is true
	if c.CaseSensitiveLike {
		accrue(sqliteparams.ConnParams.CaseSensitiveLike, c.CaseSensitiveLike)
	}
	if c.ForeignKeys {
		accrue(sqliteparams.ConnParams.ForeignKeys, c.ForeignKeys)
	}
	if c.IgnoreCheckConstraints {
		accrue(sqliteparams.ConnParams.IgnoreCheckConstraints, c.IgnoreCheckConstraints)

	}
	if c.DeferForeignKeys {
		accrue(sqliteparams.ConnParams.DeferForeignKeys, c.DeferForeignKeys)

	}
	if c.RecursiveTriggers {
		accrue(sqliteparams.ConnParams.RecursiveTriggers, c.RecursiveTriggers)
	}
	if c.WritableSchema {
		accrue(sqliteparams.ConnParams.WritableSchema, c.WritableSchema)
	}

	return buf.String()
}

func (c *Config) ToConfigParams() hypersql.ConfigParams {
	params := hypersql.ConfigParams{
		sqliteparams.ConnParams.Auth:                   cast.ToString(c.Auth),
		sqliteparams.ConnParams.AuthUser:               c.AuthUser,
		sqliteparams.ConnParams.AuthPass:               c.AuthPass,
		sqliteparams.ConnParams.AuthSalt:               c.AuthSalt,
		sqliteparams.ConnParams.AuthCrypt:              c.AuthCrypt,
		sqliteparams.ConnParams.AutoVacuum:             c.AutoVacuum,
		sqliteparams.ConnParams.BusyTimeout:            cast.ToString(c.BusyTimeout),
		sqliteparams.ConnParams.Cache:                  c.Cache,
		sqliteparams.ConnParams.CacheSize:              cast.ToString(c.CacheSize),
		sqliteparams.ConnParams.CaseSensitiveLike:      cast.ToString(c.CaseSensitiveLike),
		sqliteparams.ConnParams.DeferForeignKeys:       cast.ToString(c.DeferForeignKeys),
		sqliteparams.ConnParams.ForeignKeys:            cast.ToString(c.ForeignKeys),
		sqliteparams.ConnParams.IgnoreCheckConstraints: cast.ToString(c.IgnoreCheckConstraints),
		sqliteparams.ConnParams.Immutable:              cast.ToString(c.Immutable),
		sqliteparams.ConnParams.JournalMode:            c.JournalMode,
		sqliteparams.ConnParams.Loc:                    c.Loc,
		sqliteparams.ConnParams.LockingMode:            c.LockingMode,
		sqliteparams.ConnParams.Mode:                   c.Mode,
		sqliteparams.ConnParams.Mutex:                  c.Mutex,
		sqliteparams.ConnParams.QueryOnly:              cast.ToString(c.QueryOnly),
		sqliteparams.ConnParams.RecursiveTriggers:      cast.ToString(c.RecursiveTriggers),
		sqliteparams.ConnParams.SecureDelete:           cast.ToString(c.SecureDelete),
		sqliteparams.ConnParams.Sync:                   c.Sync,
		sqliteparams.ConnParams.TxLock:                 c.TxLock,
		sqliteparams.ConnParams.WritableSchema:         cast.ToString(c.WritableSchema),
	}
	return params
}

func ToConfigParams(c *Config) hypersql.ConfigParams {
	return c.ToConfigParams()
}

func ParseDSN(dsn string) (*Config, error) {
	pos := strings.IndexRune(dsn, '?')

	var name string
	params := make(hypersql.ConfigParams)
	if pos >= 1 {
		if query, err := url.ParseQuery(dsn[pos+1:]); err != nil {
			return nil, err
		} else {
			for k := range query {
				params[k] = query.Get(k)
			}
		}

		name = dsn[:pos]
	}

	cc := &Config{
		Name: name,
	}

	if pos == 0 {
		return cc, nil
	}

	if err := cc.HandleParams(params); err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *Config) HandleParams(params hypersql.ConfigParams) error {
	if c == nil || len(params) == 0 {
		return nil
	}

	// Options
	//var loc *time.Location
	auth := false
	authUser := ""
	authPass := ""
	authCrypt := ""
	authSalt := ""

	mutex := ""
	txlock := ""
	mode := ""

	// PRAGMA's
	autoVacuum := ""
	busyTimeout := 5000
	cache := ""
	caseSensitiveLike := false
	deferForeignKeys := false
	foreignKeys := false
	ignoreCheckConstraints := false
	journalMode := ""
	immutable := false
	lockingMode := ""
	queryOnly := false
	recursiveTriggers := false
	secureDelete := ""
	sync := ""
	writableSchema := false
	loc := ""
	cacheSize := 2000

	// Check duplicate params
	synonymParams := [][]string{
		{sqliteparams.ConnParams.AutoVacuum, sqliteparams.ConnParams.Vacuum},
		{sqliteparams.ConnParams.DeferForeignKeys, sqliteparams.ConnParams.DeferFK},
		{sqliteparams.ConnParams.CaseSensitiveLike, sqliteparams.ConnParams.CSLike},
		{sqliteparams.ConnParams.ForeignKeys, sqliteparams.ConnParams.FK},
		{sqliteparams.ConnParams.BusyTimeout, sqliteparams.ConnParams.Timeout},
		{sqliteparams.ConnParams.JournalMode, sqliteparams.ConnParams.Journal},
		{sqliteparams.ConnParams.LockingMode, sqliteparams.ConnParams.Locking},
		{sqliteparams.ConnParams.RecursiveTriggers, sqliteparams.ConnParams.RT},
		{sqliteparams.ConnParams.Synchronous, sqliteparams.ConnParams.Sync},
	}
	for _, ss := range synonymParams {
		if err := ifSynonym(params, func(keys ...string) error {
			return ErrSynonymousParam
		}, ss...); err != nil {
			return err
		}
	}

	var ifNotEmpty = func(key string, then func(value string)) {
		if val, ok := params[key]; ok && len(val) > 0 {
			then(val)
		}
	}

	// Authentication
	if _, ok := params[sqliteparams.ConnParams.Auth]; ok {
		auth = true
	}

	ifNotEmpty(sqliteparams.ConnParams.AuthUser, func(val string) {
		authUser = val
	})

	ifNotEmpty(sqliteparams.ConnParams.AuthPass, func(val string) {
		authPass = val
	})

	ifNotEmpty(sqliteparams.ConnParams.AuthCrypt, func(val string) {
		authCrypt = val
	})

	ifNotEmpty(sqliteparams.ConnParams.AuthSalt, func(val string) {
		authSalt = val
	})

	ifNotEmpty(sqliteparams.ConnParams.BusyTimeout, func(val string) {
		bt := cast.ToInt(val)
		if bt > 0 {
			busyTimeout = bt
		}
	})

	ifNotEmpty(sqliteparams.ConnParams.Loc, func(val string) {
		loc = strings.ToLower(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.DeferForeignKeys, func(val string) {
		deferForeignKeys = IsTrue(val)
	})
	ifNotEmpty(sqliteparams.ConnParams.DeferFK, func(val string) {
		deferForeignKeys = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.ForeignKeys, func(val string) {
		foreignKeys = IsTrue(val)
	})
	ifNotEmpty(sqliteparams.ConnParams.FK, func(val string) {
		foreignKeys = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.Immutable, func(val string) {
		immutable = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.QueryOnly, func(val string) {
		queryOnly = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.RecursiveTriggers, func(val string) {
		recursiveTriggers = IsTrue(val)
	})
	ifNotEmpty(sqliteparams.ConnParams.RT, func(val string) {
		recursiveTriggers = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.CaseSensitiveLike, func(val string) {
		caseSensitiveLike = IsTrue(val)
	})
	ifNotEmpty(sqliteparams.ConnParams.CSLike, func(val string) {
		caseSensitiveLike = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.IgnoreCheckConstraints, func(val string) {
		ignoreCheckConstraints = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.WritableSchema, func(val string) {
		writableSchema = IsTrue(val)
	})

	ifNotEmpty(sqliteparams.ConnParams.CacheSize, func(val string) {
		cs := cast.ToInt(val)
		if cs > 0 {
			cacheSize = cs
		}
	})

	var validVal = func(v string, ss ...string) bool {
		return slices.ContainsFunc(ss, func(e string) bool {
			return strings.EqualFold(v, e)
		})
	}

	syncHandler := func(val string) error {
		switch {
		case validVal(val,
			"0", "1", "2", "3",
			SyncFull,
			SyncNormal,
			SyncExtra,
			SyncOff):
			sync = val
			return nil
		default:
			return fmt.Errorf("invalid _sync: %v", val)
		}
	}

	autoVacuumHandler := func(val string) error {
		switch {
		case validVal(val,
			"0", "1", "2",
			AutoVacuumNone,
			AutoVacuumIncremental,
			AutoVacuumFull):
			autoVacuum = val
			return nil
		default:
			return fmt.Errorf("invalid _auto_vacuum: %v, expecting value of '0 NONE 1 FULL 2 INCREMENTAL'", val)
		}
	}

	journalHandler := func(val string) error {
		switch {
		case validVal(val,
			JournalDelete,
			JournalTruncate,
			JournalPersist,
			JournalMemory,
			JournalWal,
			JournalOff):
			journalMode = val
			return nil
		default:
			return fmt.Errorf("invalid _journal_mode: %v", val)
		}
	}

	lockingHandler := func(val string) error {
		switch {
		case validVal(val, LockingModeExclusive, LockingModeNormal):
			lockingMode = val
			return nil
		default:
			return fmt.Errorf("invalid _locking_mode: %v", val)
		}
	}

	secureDeleteHandler := func(val string) error {
		switch {
		case validVal(val,
			BoolTrueOn,
			BoolTrueYes,
			BoolTrue1,
			BoolTrueTrue):
			secureDelete = BoolTrueTrue
			return nil
		case validVal(val,
			BoolFalseOff,
			BoolFalse0,
			BoolFalseFalse,
			BoolFalseNo):
			secureDelete = BoolFalseFalse
			return nil
		case validVal(val, SecureDeleteFast):
			secureDelete = val
			return nil
		default:
			return fmt.Errorf("invalid _secure_delete: %v", val)
		}
	}

	paramHandlers := map[string]func(string) error{
		sqliteparams.ConnParams.Mode: func(val string) error {
			switch {
			case validVal(val,
				ModeRO,
				ModeRW,
				ModeRWC,
				ModeMemory):
				mode = val
				return nil
			default:
				return fmt.Errorf("invalid mode: %v", val)
			}
		},
		sqliteparams.ConnParams.Cache: func(val string) error {
			switch {
			case validVal(val, CacheShared, CachePrivate):
				cache = val
				return nil
			default:
				return fmt.Errorf("invalid cache: %v", val)
			}
		},
		sqliteparams.ConnParams.Mutex: func(val string) error {
			switch {
			case validVal(val, MutexNo, MutexFull):
				mutex = val
				return nil
			default:
				return fmt.Errorf("invalid _mutex: %v", val)
			}
		},
		sqliteparams.ConnParams.TxLock: func(val string) error {
			switch {
			case validVal(val,
				TxLockExclusive,
				TxLockDeferred,
				TxLockImmediate):
				txlock = val
				return nil
			default:
				return fmt.Errorf("invalid _txlock: %v", val)
			}
		},
		sqliteparams.ConnParams.Synchronous:  syncHandler,
		sqliteparams.ConnParams.Sync:         syncHandler,
		sqliteparams.ConnParams.AutoVacuum:   autoVacuumHandler,
		sqliteparams.ConnParams.Vacuum:       autoVacuumHandler,
		sqliteparams.ConnParams.JournalMode:  journalHandler,
		sqliteparams.ConnParams.Journal:      journalHandler,
		sqliteparams.ConnParams.SecureDelete: secureDeleteHandler,
		sqliteparams.ConnParams.LockingMode:  lockingHandler,
		sqliteparams.ConnParams.Locking:      lockingHandler,
	}
	var errorOnNotEmpty = func(key string, then func(value string) error) error {
		if v, ok := params[key]; ok && len(v) > 0 {
			return then(v)
		}
		return nil
	}
	for k, v := range paramHandlers {
		if err := errorOnNotEmpty(k, v); err != nil {
			return err
		}
	}

	c.Auth = auth
	c.AuthUser = authUser
	c.AuthPass = authPass
	c.AuthCrypt = authCrypt
	c.AuthSalt = authSalt
	c.AutoVacuum = autoVacuum
	c.BusyTimeout = busyTimeout
	c.Cache = cache
	c.CaseSensitiveLike = caseSensitiveLike
	c.DeferForeignKeys = deferForeignKeys
	c.ForeignKeys = foreignKeys
	c.TxLock = txlock
	c.Mode = mode
	c.Mutex = mutex
	c.JournalMode = journalMode
	c.Immutable = immutable
	c.LockingMode = lockingMode
	c.QueryOnly = queryOnly
	c.IgnoreCheckConstraints = ignoreCheckConstraints
	c.Sync = sync
	c.RecursiveTriggers = recursiveTriggers
	c.SecureDelete = secureDelete
	c.Loc = loc
	c.WritableSchema = writableSchema
	c.CacheSize = cacheSize

	return nil
}

func IsTrue(v any) bool {
	switch s := strings.ToLower(cast.ToString(v)); s {
	case BoolTrueOn, BoolTrueYes, BoolTrue1, BoolTrueTrue:
		return true
	default:
		return false
	}
}

func IsFalse(v any) bool {
	switch s := strings.ToLower(cast.ToString(v)); s {
	case BoolFalseOff, BoolFalseNo, BoolFalse0, BoolFalseFalse:
		return true
	default:
		return false
	}
}
