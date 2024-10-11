package hypersql

import "C"
import (
	"fmt"
	"strings"

	"github.com/blink-io/hypersql/sqlite"
	sqliteparams "github.com/blink-io/hypersql/sqlite/params"
	"github.com/spf13/cast"
)

var compatibleSQLiteDialects = []string{
	DialectSQLite,
	"sqlite3",
}

func init() {
	dialect := DialectSQLite
	//drivers[dn] = GetSQLiteDriver
	//dsners[dn] = GetSQLiteDSN
	connectors[dialect] = GetSQLiteConnector

	dialecters[dialect] = IsCompatibleSQLiteDialect
}

func (c *Config) ToSQLite() {
	c.Dialect = DialectSQLite
}

func ToSQLiteConfig(c *Config) (*sqlite.Config, error) {
	params := c.Params

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
	dupParams := [][]string{
		{sqliteparams.AutoVacuum, sqliteparams.Vacuum},
		{sqliteparams.DeferForeignKeys, sqliteparams.DeferFK},
		{sqliteparams.CaseSensitiveLike, sqliteparams.CSLike},
		{sqliteparams.ForeignKeys, sqliteparams.FK},
		{sqliteparams.BusyTimeout, sqliteparams.Timeout},
		{sqliteparams.JournalMode, sqliteparams.Journal},
		{sqliteparams.LockingMode, sqliteparams.Locking},
		{sqliteparams.RecursiveTriggers, sqliteparams.RT},
		{sqliteparams.Synchronous, sqliteparams.Sync},
	}
	for _, dps := range dupParams {
		if err := sqlite.CheckDupParam(params, dps...); err != nil {
			return nil, err
		}
	}

	// Authentication
	auth = params.Exists(sqliteparams.Auth)

	params.IfNotEmpty(sqliteparams.AuthUser, func(val string) {
		authUser = val
	})

	params.IfNotEmpty(sqliteparams.AuthPass, func(val string) {
		authPass = val
	})

	params.IfNotEmpty(sqliteparams.AuthCrypt, func(val string) {
		authCrypt = val
	})

	params.IfNotEmpty(sqliteparams.AuthSalt, func(val string) {
		authSalt = val
	})

	params.IfNotEmpty(sqliteparams.Loc, func(val string) {
		loc = strings.ToLower(val)
	})

	params.IfNotEmpty(sqliteparams.DeferForeignKeys, func(val string) {
		deferForeignKeys = sqlite.IsTrue(val)
	})
	params.IfNotEmpty(sqliteparams.DeferFK, func(val string) {
		deferForeignKeys = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.ForeignKeys, func(val string) {
		foreignKeys = sqlite.IsTrue(val)
	})
	params.IfNotEmpty(sqliteparams.FK, func(val string) {
		foreignKeys = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.Immutable, func(val string) {
		immutable = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.QueryOnly, func(val string) {
		queryOnly = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.RecursiveTriggers, func(val string) {
		recursiveTriggers = sqlite.IsTrue(val)
	})
	params.IfNotEmpty(sqliteparams.RT, func(val string) {
		recursiveTriggers = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.CaseSensitiveLike, func(val string) {
		caseSensitiveLike = sqlite.IsTrue(val)
	})
	params.IfNotEmpty(sqliteparams.CSLike, func(val string) {
		caseSensitiveLike = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.IgnoreCheckConstraints, func(val string) {
		ignoreCheckConstraints = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.WritableSchema, func(val string) {
		writableSchema = sqlite.IsTrue(val)
	})

	params.IfNotEmpty(sqliteparams.CacheSize, func(val string) {
		cs := cast.ToInt(val)
		if cs > 0 {
			cacheSize = cs
		}
	})

	syncHandler := func(val string) error {
		switch val = strings.ToLower(val); val {
		case "0", "1", "2", "3",
			sqlite.SyncFull,
			sqlite.SyncNormal,
			sqlite.SyncExtra,
			sqlite.SyncOff:
			sync = val
			return nil
		default:
			return fmt.Errorf("invalid _sync: %v", val)
		}
	}

	autoVacuumHandler := func(val string) error {
		switch val = strings.ToLower(val); val {
		case "0", "1", "2",
			sqlite.AutoVacuumNone,
			sqlite.AutoVacuumIncremental,
			sqlite.AutoVacuumFull:
			autoVacuum = val
			return nil
		default:
			return fmt.Errorf("invalid _auto_vacuum: %v, expecting value of '0 NONE 1 FULL 2 INCREMENTAL'", val)
		}
	}

	journalHandler := func(val string) error {
		switch val = strings.ToLower(val); val {
		case sqlite.JournalDelete,
			sqlite.JournalTruncate,
			sqlite.JournalMemory,
			sqlite.JournalWal,
			sqlite.JournalOff:
			journalMode = val
			return nil
		default:
			return fmt.Errorf("invalid _journal_mode: %v", val)
		}
	}

	secureDeleteHandler := func(val string) error {
		switch val = strings.ToLower(val); val {
		case sqlite.BoolTrueOn,
			sqlite.BoolTrueYes,
			sqlite.BoolTrue1,
			sqlite.BoolTrueTrue:
			secureDelete = sqlite.BoolTrueTrue
			return nil
		case sqlite.BoolFalseOff,
			sqlite.BoolFalse0,
			sqlite.BoolFalseFalse,
			sqlite.BoolFalseNo:
			secureDelete = sqlite.BoolFalseFalse
			return nil
		case sqlite.SecureDeleteFast:
			secureDelete = val
			return nil
		default:
			return fmt.Errorf("invalid _secure_delete: %v", val)
		}
	}

	paramHandlers := map[string]func(string) error{
		sqliteparams.Mode: func(val string) error {
			switch val = strings.ToLower(val); val {
			case sqlite.ModeRO,
				sqlite.ModeRW,
				sqlite.ModeRWC,
				sqlite.ModeMemory:
				mode = val
				return nil
			default:
				return fmt.Errorf("invalid mode: %v", val)
			}
		},
		sqliteparams.Cache: func(val string) error {
			switch val = strings.ToLower(val); val {
			case sqlite.CacheShared, sqlite.CachePrivate:
				cache = val
				return nil
			default:
				return fmt.Errorf("invalid cache: %v", val)
			}
		},
		sqliteparams.Mutex: func(val string) error {
			switch val = strings.ToLower(val); val {
			case sqlite.MutexNo, sqlite.MutexFull:
				mutex = val
				return nil
			default:
				return fmt.Errorf("invalid _mutex: %v", val)
			}
		},
		sqliteparams.TxLock: func(val string) error {
			switch val = strings.ToLower(val); val {
			case sqlite.TxLockExclusive,
				sqlite.TxLockDeferred,
				sqlite.TxLockImmediate:
				txlock = val
				return nil
			default:
				return fmt.Errorf("invalid _txlock: %v", val)
			}
		},
		sqliteparams.Synchronous:  syncHandler,
		sqliteparams.Sync:         syncHandler,
		sqliteparams.AutoVacuum:   autoVacuumHandler,
		sqliteparams.Vacuum:       autoVacuumHandler,
		sqliteparams.JournalMode:  journalHandler,
		sqliteparams.Journal:      journalHandler,
		sqliteparams.SecureDelete: secureDeleteHandler,
	}
	for k, v := range paramHandlers {
		if err := params.IfNotEmptyWithError(k, v); err != nil {
			return nil, err
		}
	}

	if len(c.User) > 0 {
		authUser = c.User
	}

	if len(c.Password) > 0 {
		authPass = c.Password
	}

	cc := &sqlite.Config{
		Name:                   c.Name,
		Auth:                   auth,
		AuthUser:               authUser,
		AuthPass:               authPass,
		AuthCrypt:              authCrypt,
		AuthSalt:               authSalt,
		AutoVacuum:             autoVacuum,
		BusyTimeout:            busyTimeout,
		Cache:                  cache,
		CaseSensitiveLike:      caseSensitiveLike,
		DeferForeignKeys:       deferForeignKeys,
		ForeignKeys:            foreignKeys,
		TxLock:                 txlock,
		Mode:                   mode,
		Mutex:                  mutex,
		JournalMode:            journalMode,
		Immutable:              immutable,
		LockingMode:            lockingMode,
		QueryOnly:              queryOnly,
		IgnoreCheckConstraints: ignoreCheckConstraints,
		Sync:                   sync,
		RecursiveTriggers:      recursiveTriggers,
		SecureDelete:           secureDelete,
		Loc:                    loc,
		WritableSchema:         writableSchema,
		CacheSize:              cacheSize,
	}

	return cc, nil
}
