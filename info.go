package hypersql

type (
	DBInfo struct {
		Name    string
		Dialect string
	}

	WithDBInfo interface {
		DBInfo() DBInfo
	}
)

func NewDBInfo(c *Config) DBInfo {
	return DBInfo{
		Name:    c.Name,
		Dialect: c.Dialect,
	}
}
