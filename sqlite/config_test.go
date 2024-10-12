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

func TestDSN(t *testing.T) {
	c := &Config{
		Name:              "file:sqlite.demo",
		Auth:              true,
		AuthUser:          "user",
		AuthPass:          "pass",
		AuthCrypt:         CryptSHA256,
		AuthSalt:          uuid.NewString(),
		Cache:             CachePrivate,
		Mode:              ModeMemory,
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

	t.Run("Format DSN 1", func(t *testing.T) {
		c1 := &Config{
			Name: "sqlite.db",
			Auth: true,
			Mode: ModeMemory,
		}
		dsn1 := c1.FormatDSN()
		fmt.Println(dsn1)
	})

	t.Run("Parse DSN", func(t *testing.T) {
		cc, err := ParseDSN(dsn)
		require.NoError(t, err)
		require.NotNil(t, cc)
		require.Equal(t, cc, c)
	})

	t.Run("Exec Query", func(t *testing.T) {
		db, err := sql.Open("sqlite", dsn)

		require.NoError(t, err)
		require.NotNil(t, db)

		defer db.Close()

		row := db.QueryRow("select sqlite_version()")
		var ver string
		row.Scan(&ver)

		fmt.Println("SQLite version: ", ver)
	})
}
