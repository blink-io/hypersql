package sqlite

import (
	"bytes"
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

func writeParamPair(buf *bytes.Buffer, joiner, key string, value any) {
	buf.WriteString(joiner)
	buf.WriteString(key)
	buf.WriteString("=")
	buf.WriteString(cast.ToString(value))
}

func (c *Config) FormatDSN() string {
	if !strings.HasPrefix(c.Name, "file:") {
		return c.Name
	}

	joiner := "&"
	var buf bytes.Buffer
	buf.WriteString(c.Name)
	buf.WriteString("?")
	writeParamPair(&buf, "", params.ConnParams.Immutable, c.Immutable)

	if c.QueryOnly {
		writeParamPair(&buf, joiner, params.ConnParams.QueryOnly, c.QueryOnly)
	}
	if c.Auth {
		buf.WriteString("&_auth=")
		if len(c.AuthUser) > 0 {
			writeParamPair(&buf, joiner, params.ConnParams.AuthUser, c.AuthUser)
		}
		if len(c.AuthPass) > 0 {
			writeParamPair(&buf, joiner, params.ConnParams.AuthPass, c.AuthPass)
		}
		if len(c.AuthCrypt) > 0 {
			writeParamPair(&buf, joiner, params.ConnParams.AuthCrypt, c.AuthCrypt)
		}
		if len(c.AuthSalt) > 0 {
			writeParamPair(&buf, joiner, params.ConnParams.AuthSalt, c.AuthSalt)
		}
	}
	if len(c.Cache) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.Cache, c.Cache)
	}
	if len(c.Mode) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.Mode, c.Mode)
	}
	if len(c.Mutex) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.Mutex, c.Mutex)
	}
	if c.BusyTimeout > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.BusyTimeout, c.BusyTimeout)
	}
	if len(c.AutoVacuum) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.AutoVacuum, c.AutoVacuum)
	}
	if len(c.JournalMode) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.JournalMode, c.JournalMode)
	}
	if len(c.LockingMode) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.LockingMode, c.LockingMode)
	}
	if len(c.SecureDelete) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.SecureDelete, c.SecureDelete)
	}
	if len(c.Loc) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.Loc, c.Loc)
	}
	if len(c.Sync) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.Sync, c.Sync)
	}
	if len(c.TxLock) > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.TxLock, c.TxLock)
	}
	if c.CacheSize > 0 {
		writeParamPair(&buf, joiner, params.ConnParams.CacheSize, c.CacheSize)
	}

	// Append other params when it is true
	if c.CaseSensitiveLike {
		writeParamPair(&buf, joiner, params.ConnParams.CaseSensitiveLike, c.CaseSensitiveLike)
	}
	if c.ForeignKeys {
		writeParamPair(&buf, joiner, params.ConnParams.ForeignKeys, c.ForeignKeys)
	}
	if c.IgnoreCheckConstraints {
		writeParamPair(&buf, joiner, params.ConnParams.IgnoreCheckConstraints, c.IgnoreCheckConstraints)

	}
	if c.DeferForeignKeys {
		writeParamPair(&buf, joiner, params.ConnParams.DeferForeignKeys, c.DeferForeignKeys)

	}
	if c.RecursiveTriggers {
		writeParamPair(&buf, joiner, params.ConnParams.RecursiveTriggers, c.RecursiveTriggers)
	}
	if c.WritableSchema {
		writeParamPair(&buf, joiner, params.ConnParams.WritableSchema, c.WritableSchema)
	}

	return buf.String()
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
