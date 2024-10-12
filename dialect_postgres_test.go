package hypersql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDialectPostgresDSN(t *testing.T) {
	var ctx = context.Background()

	cc := &Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "postgres",
	}

	t.Run("ToDSN", func(t *testing.T) {
		dsn, err := ToPostgresDSN(ctx, cc)
		require.NoError(t, err)
		require.Equal(t, "dbname='postgres' host='localhost' password='postgres' port='5432' user='postgres'", dsn)
	})
}
