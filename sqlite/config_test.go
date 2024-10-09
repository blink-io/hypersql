package sqlite

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
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
}
