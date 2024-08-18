//go:build use_cgo && !cgo_ext

package hypersql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHypersql_SQLite_CGO_1(t *testing.T) {
	drv := getRawSQLiteDriver()
	sql.Register("sqlite_cgo_free", drv)
	db, err := sql.Open("sqlite_cgo_free", "file:demo.db?_txlock=immediate")
	require.NoError(t, err)

	r := db.QueryRow("select sqlite_version()")
	var str string
	require.NoError(t, r.Scan(&str))

	fmt.Println("SQLite Version:", str)
}
