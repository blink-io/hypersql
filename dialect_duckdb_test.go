package hypersql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDuckDB_1(t *testing.T) {
	drv := getRawDuckDBDriver()
	sql.Register("duckdb_drv", drv)
	db, err := sql.Open("duckdb_drv", "")
	require.NoError(t, err)

	r := db.QueryRow("select version() as version, current_schema() as schema")
	var version, schema string
	require.NoError(t, r.Scan(&version, &schema))

	fmt.Println("DuckDB version:", version)
	fmt.Println("DuckDB schema:", schema)
}
