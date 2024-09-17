package hypersql

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestPgxType_1(t *testing.T) {
	tmz := pgtype.Timestamptz{}
	fmt.Println(tmz)
}
