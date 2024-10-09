package sqlite

import (
	"errors"
	"strings"

	"github.com/blink-io/hypersql/sqlite/params"
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
	ErrDuplicateParam = errors.New("duplicate param")
)

// Config is the string for the SQLite database.
// See https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string
type Config struct {
	Name                   string
	Auth                   bool
	AuthUser               string
	AuthPass               string
	AuthCrypt              string
	AuthSalt               string
	Cache                  string
	Mode                   string
	Mutex                  string
	CaseSensitiveLike      bool
	Immutable              bool
	QueryOnly              bool
	ForeignKeys            bool
	IgnoreCheckConstraints bool
	DeferForeignKeys       bool
	BusyTimeout            int
	AutoVacuum             string
	JournalMode            string
	LockingMode            string
	RecursiveTriggers      bool
	SecureDelete           string
	Sync                   string
	Loc                    string
	TxLock                 string
	WritableSchema         bool
	CacheSize              int
}

func (c *Config) FormatDSN() string {
	if !strings.HasPrefix(c.Name, "file:") {
		return c.Name
	}

	var kv = func(prefix, suffix, key string, value any) string {
		return prefix + key + "=" + cast.ToString(value) + suffix
	}
	joint := "&"
	var dsnBuilder strings.Builder
	dsnBuilder.WriteString(c.Name)
	dsnBuilder.WriteString("?")
	dsnBuilder.WriteString(kv("", "", params.QueryOnly, c.QueryOnly))
	dsnBuilder.WriteString(kv("", "", params.Immutable, c.Immutable))

	if c.Auth {
		dsnBuilder.WriteString("&_auth=")
		if len(c.AuthUser) > 0 {
			dsnBuilder.WriteString(kv(joint, "", params.AuthUser, c.AuthUser))
		}
		if len(c.AuthPass) > 0 {
			dsnBuilder.WriteString(kv(joint, "", params.AuthPass, c.AuthPass))
		}
		if len(c.AuthCrypt) > 0 {
			dsnBuilder.WriteString(kv(joint, "", params.AuthCrypt, c.AuthCrypt))
		}
		if len(c.AuthSalt) > 0 {
			dsnBuilder.WriteString(kv(joint, "", params.AuthSalt, c.AuthSalt))
		}
	}
	if len(c.Cache) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.Cache, c.Cache))
	}
	if len(c.Mode) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.Mode, c.Mode))
	}
	if len(c.Mutex) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.Mutex, c.Mutex))
	}
	if c.BusyTimeout > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.BusyTimeout, c.BusyTimeout))
	}
	if len(c.AutoVacuum) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.AutoVacuum, c.AutoVacuum))
	}
	if len(c.JournalMode) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.JournalMode, c.JournalMode))
	}
	if len(c.LockingMode) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.LockingMode, c.LockingMode))
	}
	if len(c.SecureDelete) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.SecureDelete, c.SecureDelete))
	}
	if len(c.Loc) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.Loc, c.Loc))
	}
	if len(c.Sync) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.Sync, c.Sync))
	}
	if len(c.TxLock) > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.TxLock, c.TxLock))
	}
	if c.CacheSize > 0 {
		dsnBuilder.WriteString(kv(joint, "", params.CacheSize, c.CacheSize))
	}

	// Append other params when it is true
	if c.CaseSensitiveLike {
		dsnBuilder.WriteString(kv(joint, "", params.CaseSensitiveLike, c.CaseSensitiveLike))
	}
	if c.ForeignKeys {
		dsnBuilder.WriteString(kv(joint, "", params.ForeignKeys, c.ForeignKeys))
	}
	if c.IgnoreCheckConstraints {
		dsnBuilder.WriteString(kv(joint, "", params.IgnoreCheckConstraints, c.IgnoreCheckConstraints))

	}
	if c.DeferForeignKeys {
		dsnBuilder.WriteString(kv(joint, "", params.DeferForeignKeys, c.DeferForeignKeys))

	}
	if c.RecursiveTriggers {
		dsnBuilder.WriteString(kv(joint, "", params.RecursiveTriggers, c.RecursiveTriggers))

	}
	if c.WritableSchema {
		dsnBuilder.WriteString(kv(joint, "", params.WritableSchema, c.WritableSchema))
	}

	return dsnBuilder.String()
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

func dupCheck(params map[string]string, keys ...string) bool {
	if len(keys) > 1 {
		var c = 0
		for _, key := range keys {
			if params[key] != "" {
				c++
				if c > 1 {
					return true
				}
			}
		}
	}
	return false
}

func CheckDupParam(params map[string]string, keys ...string) error {
	if dupCheck(params, keys...) {
		return ErrDuplicateParam
	}
	return nil
}
