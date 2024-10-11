package hypersql

import (
	"context"
	"fmt"
	"testing"

	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptions_1(t *testing.T) {
	var o *Config
	o = SetupConfig(o)
	require.NotNil(t, o)
}

func TestOptions_Validate(t *testing.T) {
	var o *Config
	o = SetupConfig(o)
	verr := o.Validate(context.Background())
	require.Error(t, verr)
	assert.ErrorIs(t, verr, ErrUnsupportedDialect)
}

func TestParseURL_1(t *testing.T) {

	t.Run("MySQL URL", func(t *testing.T) {

	})

	t.Run("Postgres URL", func(t *testing.T) {
		urlstr := "postgresql://user:pass@localhost:5432/mydatabase?sslmode=disable&TimeZone=Asia/Shanghai"

		cc, err := ParseURL(urlstr)
		require.NoError(t, err)

		fmt.Println("DSN: ", cc.dsn)

		pgcc, err := ToPostgresConfigFromDSN(cc.dsn)
		require.NoError(t, err)
		require.NotNil(t, pgcc)

		litter.Dump(cc)
	})

	t.Run("dd", func(t *testing.T) {

	})
}
