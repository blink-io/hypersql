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
	mssqlparams "github.com/blink-io/hypersql/sqlserver/params"
	"github.com/qustavo/sqlhooks/v2/hooks/loghooks"
	"github.com/stretchr/testify/assert"
)

func TestSQLServer_Driver_1(t *testing.T) {
	c := &Config{
		Dialect:  DialectSQLServer,
		Name:     "test",
		User:     "sa",
		Host:     "localhost",
		Port:     1433,
		Password: "Heison99188",
		Params: ConfigParams{
			//sqliteparams.ConnParams.Cache: sqlite.CacheShared,
			//sqliteparams.ConnParams.Mode:  sqlite.ModeMemory,
			mssqlparams.ConnParams.Encrypt: "3",
		},
		DriverHooks: DriverHooks{
			loghooks.New(),
		},
		DriverWrappers: DriverWrappers{
			func(drv driver.Driver) driver.Driver {
				return otelsql.WrapDriver(drv)
			},
			func(drv driver.Driver) driver.Driver {
				slog.Info("=======================")
				return drv
			},
		},
	}

	t.Run("sqlserver success", func(t *testing.T) {
		c.AfterHandlers = AfterHandlers{
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
		//sql := "select @@Version"
		sql := "select * from users"
		row := db.QueryRow(sql)
		assert.NoError(t, row.Scan(&ver))

		fmt.Println("sqlserver version: ", ver)
	})

	t.Run("sqlserver fail", func(t *testing.T) {
		c.AfterHandlers = AfterHandlers{
			func(ctx context.Context, db *sql.DB) error {
				return errors.New("throw an error for fail")
			},
		}

		db, err := NewSqlDB(c)
		assert.Errorf(t, err, err.Error())
		assert.Nil(t, db)
	})
}
