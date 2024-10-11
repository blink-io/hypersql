//go:build sqlite

package hypersql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
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

	t.Run("success", func(t *testing.T) {
		c.SqlDBHandlers = SqlDBHandlers{
			func(ctx context.Context, db *sql.DB) error {
				return otelsql.RegisterDBStatsMetrics(db)
			},
			func(ctx context.Context, db *sql.DB) error {
				slog.Info("sql.DB is created successfully")
				return nil
			},
		}

		db, err := NewSqlDB(c)
		assert.NoError(t, err)

		var ver string
		sql := "select sqlite_version()"
		row := db.QueryRow(sql)
		assert.NoError(t, row.Scan(&ver))

		fmt.Println("sqlite version: ", ver)
	})

	t.Run("fail", func(t *testing.T) {
		c.SqlDBHandlers = SqlDBHandlers{
			func(ctx context.Context, db *sql.DB) error {
				return errors.New("throw an error for fail")
			},
		}

		db, err := NewSqlDB(c)
		assert.Errorf(t, err, err.Error())
		assert.Nil(t, db)
	})
}
