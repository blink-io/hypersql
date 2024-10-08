package params

// Source: https://github.com/mattn/go-sqlite3#connection-string
const (
	Auth = "_auth"

	AuthUser = "_auth_user"

	AuthPass = "_auth_pass"

	AuthCrypt = "_auth_crypt"

	AuthSalt = "_auth_salt"

	AutoVacuum = "_auto_vacuum"

	Vacuum = "_vacuum"

	BusyTimeout = "_busy_timeout"

	Timeout = "_timeout"

	CaseSensitiveLike = "_case_sensitive_like"

	CSLike = "_cslike"

	DeferForeignKeys = "_defer_foreign_keys" // boolean

	DeferFK = "_defer_fk" // boolean

	ForeignKeys = "_foreign_keys" // boolean

	FK = "_fk" // boolean

	IgnoreCheckConstraints = "_ignore_check_constraints" // boolean

	JournalMode = "_journal_mode"

	Journal = "_journal"

	Immutable = "immutable"

	LockingMode = "_locking_mode"

	Locking = "_locking"

	Loc = "_loc"

	Mode = "mode"

	Mutex = "_mutex"

	QueryOnly = "_query_only" // boolean

	RecursiveTriggers = "_recursive_triggers" // boolean

	RT = "_rt" // boolean

	SecureDelete = "_secure_delete"

	SharedCacheMode = "cache"

	Synchronous = "_synchronous"

	Sync = "_sync"

	TxLock = "_txlock"

	WritableSchema = "_writable_schema"

	CacheSize = "_cache_size"

	Cache = "cache"
)
