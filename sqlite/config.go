package sqlite

// Config is the configuration for the SQLite database.
// See https://pkg.go.dev/github.com/mattn/go-sqlite3#readme-connection-string
type Config struct {
	Name                   string
	User                   string
	Password               string
	Crypt                  string
	Salt                   string
	Mode                   string
	Mutex                  string
	CaseSensitiveLike      bool
	Immutable              bool
	Location               string
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
	TxLock                 string
	WritableSchema         bool
	CacheSize              int
}

func (c *Config) FormatDSN() string {
	return c.Name
}
