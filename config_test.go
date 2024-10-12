package hypersql

import (
	"context"
	"fmt"
	"net"
	nurl "net/url"
	"sort"
	"strings"
	"testing"

	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xo/dburl"
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
	var ctx = context.Background()

	t.Run("MySQL URL", func(t *testing.T) {

	})

	t.Run("Postgres URL", func(t *testing.T) {
		urlstr := "postgresql://user:pass@localhost:5432/mydatabase?sslmode=disable&TimeZone=Asia/Shanghai"

		cc, err := ParseURL(urlstr)
		require.NoError(t, err)

		dsn, err := ToPostgresDSN(ctx, cc)
		require.NoError(t, err)

		fmt.Println("DSN: ", dsn)

		pgcc, err := ToPostgresConfigFromDSN(dsn)
		require.NoError(t, err)
		require.NotNil(t, pgcc)

		litter.Dump(cc)
	})

	t.Run("SQLite URL Parse", func(t *testing.T) {
		url := "sqlite:/path/to/file.db"
		uu, err := dburl.Parse(url)
		require.NoError(t, err)
		require.NotNil(t, uu)

	})

	t.Run("SQLite URL Parse 2", func(t *testing.T) {
		url := "file:myfile.sqlite3?loc=auto"
		uu, err := dburl.Parse(url)
		require.NoError(t, err)
		require.NotNil(t, uu)
	})

	t.Run("dd", func(t *testing.T) {
		urlstr := "postgresql://user:pass@localhost:5432/mydatabase?sslmode=disable&TimeZone=Asia/Shanghai"
		ParseURL := func(url string) (string, error) {
			u, err := nurl.Parse(url)
			if err != nil {
				return "", err
			}

			if u.Scheme != "postgres" && u.Scheme != "postgresql" {
				return "", fmt.Errorf("invalid connection protocol: %s", u.Scheme)
			}

			var kvs []string
			escaper := strings.NewReplacer(`'`, `\'`, `\`, `\\`)
			accrue := func(k, v string) {
				if v != "" {
					kvs = append(kvs, k+"='"+escaper.Replace(v)+"'")
				}
			}

			if u.User != nil {
				v := u.User.Username()
				accrue("user", v)

				v, _ = u.User.Password()
				accrue("password", v)
			}

			if host, port, err := net.SplitHostPort(u.Host); err != nil {
				accrue("host", u.Host)
			} else {
				accrue("host", host)
				accrue("port", port)
			}

			if u.Path != "" {
				accrue("dbname", u.Path[1:])
			}

			q := u.Query()
			for k := range q {
				accrue(k, q.Get(k))
			}

			sort.Strings(kvs) // Makes testing easier (not a performance concern)
			return strings.Join(kvs, " "), nil
		}

		dsn, err := ParseURL(urlstr)
		require.NoError(t, err)

		fmt.Println("DSN: ", dsn)
	})
}
