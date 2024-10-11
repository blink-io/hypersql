package hypersql

import (
	"maps"

	"github.com/spf13/cast"
	"github.com/xo/dburl"
)

func ParseURL(url string) (*Config, error) {
	uu, err := dburl.Parse(url)
	if err != nil {
		return nil, err
	}

	dialect := GetFormalDialect(uu.Scheme)
	if len(dialect) == 0 {
		return nil, ErrUnsupportedDialect
	}

	c := &Config{
		Transport: uu.Transport,
		Dialect:   dialect,
		Name:      uu.EscapedPath(),
		Host:      uu.Host,
		Port:      cast.ToInt(uu.Port),
		Params:    make(ConfigParams),
		url:       url,
		dsn:       uu.DSN,
	}

	if ui := uu.User; ui != nil {
		c.User = ui.Username()
		pwd, _ := ui.Password()
		c.Password = pwd
	}

	query := uu.Query()
	for k := range maps.Keys(query) {
		c.Params[k] = query.Get(k)
	}

	return c, nil
}
