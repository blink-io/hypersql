//go:build sqlite

package sqlite

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSqlite_DSN_1(t *testing.T) {
	c := &Config{
		Name:              "file:sqlite.demo",
		Auth:              true,
		AuthUser:          "user",
		AuthPass:          "pass",
		AuthCrypt:         CryptSHA256,
		AuthSalt:          uuid.NewString(),
		Cache:             CachePrivate,
		Mode:              ModeRWC,
		Mutex:             MutexFull,
		CaseSensitiveLike: true,
		Immutable:         true,
		QueryOnly:         true,
		AutoVacuum:        AutoVacuumIncremental,
		Sync:              SyncNormal,
		TxLock:            TxLockExclusive,
		JournalMode:       JournalMemory,
		LockingMode:       LockingModeExclusive,
		SecureDelete:      SecureDeleteFast,
		Loc:               LocAuto,
		CacheSize:         5000,
		BusyTimeout:       3000,
	}
	dsn := c.FormatDSN()
	fmt.Println(dsn)

	db, err := sql.Open("sqlite", dsn)

	require.NoError(t, err)
	require.NotNil(t, db)

	defer db.Close()

	row := db.QueryRow("select sqlite_version()")
	var ver string
	row.Scan(&ver)

	fmt.Println("Sqlite version: ", ver)
}
