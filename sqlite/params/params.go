package params

import (
	"strings"
)

// ConnParams represents the supported connection parameters.
// Source: https://github.com/mattn/go-sqlite3#connection-string
var ConnParams = connParams{

	// https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-user-authentication
	Auth: "_auth",

	AuthUser: "_auth_user",

	AuthPass: "_auth_pass",

	AuthCrypt: "_auth_crypt",

	AuthSalt: "_auth_salt",

	AutoVacuum: "_auto_vacuum",

	Vacuum: "_vacuum",

	BusyTimeout: "_busy_timeout",

	Timeout: "_timeout",

	CaseSensitiveLike: "_case_sensitive_like",

	CSLike: "_cslike",

	DeferForeignKeys: "_defer_foreign_keys", // boolean

	DeferFK: "_defer_fk", // boolean

	ForeignKeys: "_foreign_keys", // boolean

	FK: "_fk", // boolean

	IgnoreCheckConstraints: "_ignore_check_constraints", // boolean

	JournalMode: "_journal_mode",

	Journal: "_journal",

	Immutable: "immutable",

	LockingMode: "_locking_mode",

	Locking: "_locking",

	Loc: "_loc",

	Mode: "mode",

	Cache: "cache",

	Mutex: "_mutex",

	QueryOnly: "_query_only", // boolean

	RecursiveTriggers: "_recursive_triggers", // boolean

	RT: "_rt", // boolean

	SecureDelete: "_secure_delete",

	SharedCacheMode: "cache",

	Synchronous: "_synchronous",

	Sync: "_sync",

	TxLock: "_txlock",

	WritableSchema: "_writable_schema",

	CacheSize: "_cache_size",
}

type connParams struct {
	Auth                   string
	AuthUser               string
	AuthPass               string
	AuthCrypt              string
	AuthSalt               string
	AutoVacuum             string
	Vacuum                 string
	BusyTimeout            string
	Timeout                string
	CaseSensitiveLike      string
	CSLike                 string
	DeferForeignKeys       string
	DeferFK                string
	ForeignKeys            string
	FK                     string
	IgnoreCheckConstraints string
	JournalMode            string
	Journal                string
	Immutable              string
	Loc                    string
	Locale                 string
	LockingMode            string
	Locking                string
	Mode                   string
	Cache                  string
	Mutex                  string
	QueryOnly              string
	RecursiveTriggers      string
	RT                     string
	SecureDelete           string
	SharedCacheMode        string
	Synchronous            string
	Sync                   string
	TxLock                 string
	WritableSchema         string
	CacheSize              string
}

func (p connParams) Exists(key string) bool {
	switch key := strings.ToLower(key); key {
	case p.Auth,
		p.AuthUser,
		p.AuthPass,
		p.AuthCrypt,
		p.AuthSalt,
		p.AutoVacuum,
		p.Vacuum,
		p.BusyTimeout,
		p.Timeout,
		p.CaseSensitiveLike,
		p.CSLike,
		p.DeferForeignKeys,
		p.DeferFK,
		p.ForeignKeys,
		p.FK,
		p.IgnoreCheckConstraints,
		p.JournalMode,
		p.Journal,
		p.Immutable,
		p.Loc,
		p.Locale,
		p.LockingMode,
		p.Locking,
		p.Mode,
		p.Cache,
		p.Mutex,
		p.QueryOnly,
		p.RecursiveTriggers,
		p.RT,
		p.SecureDelete,
		p.SharedCacheMode,
		p.Synchronous,
		p.Sync,
		p.TxLock,
		p.WritableSchema,
		p.CacheSize:
		return true
	default:
		return false
	}
}
