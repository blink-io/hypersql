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

	if val := params.Get(sqliteparams.AuthUser); val != "" {
		authUser = val
	}
	if val := params.Get(sqliteparams.AuthPass); val != "" {
		authPass = val
	}
	if val := params.Get(sqliteparams.AuthCrypt); val != "" {
		authCrypt = val
	}
	if val := params.Get(sqliteparams.AuthSalt); val != "" {
		authSalt = val
	}

	// mode
	if val := params.Get(sqliteparams.Mode); val != "" {
		switch val = strings.ToLower(val); val {
		case sqlite.ModeRO,
			sqlite.ModeRW,
			sqlite.ModeRWC,
			sqlite.ModeMemory:
			mode = val
		default:
			return nil, fmt.Errorf("invalid mode: %v", val)
		}
	}

	// _loc
	if val := params.Get(sqliteparams.Loc); val != "" {
		loc = strings.ToLower(val)
	}

	// cache
	if val := params.Get(sqliteparams.Cache); val != "" {
		switch val = strings.ToLower(val); val {
		case sqlite.CacheShared, sqlite.CachePrivate:
			cache = val
		default:
			return nil, fmt.Errorf("invalid cache: %v", val)
		}
	}

	// _mutex
	if val := params.Get(sqliteparams.Mutex); val != "" {
		switch val = strings.ToLower(val); val {
		case sqlite.MutexNo, sqlite.MutexFull:
			mutex = val
		default:
			return nil, fmt.Errorf("invalid _mutex: %v", val)
		}
	}

	// _txlock
	if val := params.Get(sqliteparams.TxLock); val != "" {
		switch val = strings.ToLower(val); val {
		case sqlite.TxLockExclusive,
			sqlite.TxLockDeferred,
			sqlite.TxLockImmediate:
			txlock = val
		default:
			return nil, fmt.Errorf("invalid _txlock: %v", val)
		}
	}

	// _sync
	if val := params.Get(sqliteparams.Synchronous, sqliteparams.Sync); val != "" {
		switch val = strings.ToLower(val); val {
		case "0", "1", "2", "3",
			sqlite.SyncFull,
			sqlite.SyncNormal,
			sqlite.SyncExtra,
			sqlite.SyncOff:
			sync = val
		default:
			return nil, fmt.Errorf("invalid _sync: %v", val)
		}
	}

	if val := params.Get(sqliteparams.AutoVacuum, sqliteparams.Vacuum); val != "" {
		switch val = strings.ToLower(val); val {
		case "0", "1", "2",
			sqlite.AutoVacuumNone,
			sqlite.AutoVacuumIncremental,
			sqlite.AutoVacuumFull:
			autoVacuum = val
		default:
			return nil, fmt.Errorf("invalid _auto_vacuum: %v, expecting value of '0 NONE 1 FULL 2 INCREMENTAL'", val)
		}
	}

	if val := params.Get(sqliteparams.JournalMode, sqliteparams.Journal); val != "" {
		switch val = strings.ToLower(val); val {
		case sqlite.JournalDelete,
			sqlite.JournalTruncate,
			sqlite.JournalMemory,
			sqlite.JournalWal,
			sqlite.JournalOff:
			journalMode = val
		default:
			return nil, fmt.Errorf("invalid _journal_mode: %v", val)
		}
	}

	if val := params.Get(sqliteparams.CaseSensitiveLike, sqliteparams.CSLike); val != "" {
		caseSensitiveLike = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.DeferForeignKeys, sqliteparams.DeferFK); val != "" {
		deferForeignKeys = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.ForeignKeys, sqliteparams.FK); val != "" {
		foreignKeys = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.Immutable); val != "" {
		immutable = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.QueryOnly); val != "" {
		queryOnly = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.RecursiveTriggers, sqliteparams.RT); val != "" {
		recursiveTriggers = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.IgnoreCheckConstraints); val != "" {
		ignoreCheckConstraints = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.WritableSchema); val != "" {
		writableSchema = sqlite.IsTrue(val)
	}

	if val := params.Get(sqliteparams.CacheSize); val != "" {
		cs := cast.ToInt(val)
		if cs > 0 {
			cacheSize = cs
		}
	}

	if val := params.Get(sqliteparams.SecureDelete); val != "" {
		switch val = strings.ToLower(val); val {
		case sqlite.BoolTrueOn,
			sqlite.BoolTrueYes,
			sqlite.BoolTrue1,
			sqlite.BoolTrueTrue:
			secureDelete = sqlite.BoolTrueTrue
		case sqlite.BoolFalseOff,
			sqlite.BoolFalse0,
			sqlite.BoolFalseFalse,
			sqlite.BoolFalseNo:
			secureDelete = sqlite.BoolFalseFalse
		case sqlite.SecureDeleteFast:
			secureDelete = val
		default:
			return nil, fmt.Errorf("invalid _secure_delete: %v", val)
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
