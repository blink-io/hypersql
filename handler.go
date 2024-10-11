package hypersql

import (
	"context"
	"database/sql"
)

type (
	SqlDBHandler func(context.Context, *sql.DB) error

	SqlDBHandlers []SqlDBHandler
)
