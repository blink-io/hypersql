package hypersql

import (
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEqualFold(t *testing.T) {
	var str1 = "Hello"
	var str2 = "HELLo"
	var flag = strings.EqualFold(str1, str2)

	assert.True(t, flag)
}

func TestLocation_1(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)
	require.NotNil(t, loc)
}

func TestPgxParseConfig(t *testing.T) {
	str := "sslmode=disable"
	cc, err := pgx.ParseConfig(str)
	require.NoError(t, err)
	require.NotNil(t, cc)
}

func TestCastDuration(t *testing.T) {
	d1 := cast.ToDuration("1m4s")
	d2 := 64 * time.Second
	assert.Equal(t, d1, d2)
}
