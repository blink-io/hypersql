package sqlserver

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/msdsn"
	"github.com/stretchr/testify/require"
)

func TestDSN(t *testing.T) {
	c := msdsn.Config{
		Database:    "db",
		Instance:    "inst1",
		Host:        "localhost",
		Port:        1334,
		User:        "user",
		Password:    "pass",
		AppName:     "app_name",
		Workstation: "ws_",
		ServerSPN:   "srv_spn",
	}

	urlstr := c.URL().String()

	fmt.Println("url: ", urlstr)

	cc, err := msdsn.Parse(urlstr)
	require.NoError(t, err)
	require.NotNil(t, cc)
}

func TestSQLServer_Connect(t *testing.T) {
	sqlstr := "select @@Version;"
	dsn := "sqlserver://sa:Heison99188@localhost:1433?database=master&connection+timeout=30"
	db, err := sql.Open("sqlserver", dsn)
	require.NoError(t, err)

	defer db.Close()

	var ver string
	err = db.QueryRow(sqlstr).Scan(&ver)
	require.NoError(t, err)
}
