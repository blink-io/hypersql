//go:build sqlite

package hypersql

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/XSAM/otelsql"
	"github.com/blink-io/hypersql/sqlite"
	sqliteparams "github.com/blink-io/hypersql/sqlite/params"
	"github.com/qustavo/sqlhooks/v2/hooks/loghooks"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestSqlite_Driver_1(t *testing.T) {
	c := &Config{
		Dialect: DialectSQLite,
		Name:    "file:sqlite.db",
		Params: ConfigParams{
			sqliteparams.Cache: sqlite.CacheShared,
			sqliteparams.Mode:  sqlite.ModeMemory,
		},
		DriverHooks: DriverHooks{
			loghooks.New(),
		},
		DriverWrappers: DriverWrappers{
			func(drv driver.Driver) driver.Driver {
				return otelsql.WrapDriver(drv)
			},
		},
	}

	db, err := NewSqlDB(c)
	assert.NoError(t, err)

	var ver string
	sql := "select sqlite_version()"
	row := db.QueryRow(sql)
	assert.NoError(t, row.Scan(&ver))

	fmt.Println("sqlite version: ", ver)
}
