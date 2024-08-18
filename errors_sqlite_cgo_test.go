//go:build use_cgo

package hypersql

import (
	"fmt"
	"testing"

	"github.com/mattn/go-sqlite3"
	"github.com/spf13/cast"
)

func TestSQLite3ErrNo_Actual(t *testing.T) {
	nn := map[string]sqlite3.ErrNoExtended{
		"notnull": sqlite3.ErrConstraintNotNull,
		"unique":  sqlite3.ErrConstraintUnique,
		"check":   sqlite3.ErrConstraintCheck,
		"fk":      sqlite3.ErrConstraintForeignKey,
		"pk":      sqlite3.ErrConstraintPrimaryKey,
	}

	for k, v := range nn {
		vi := int(v)
		fmt.Println("Name ---> ", k, "  Num ---> ", vi, "Num in String ---> ", cast.ToString(vi))
	}
}
