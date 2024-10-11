package hypersql

import (
	"context"
	"database/sql"
)

type (
	AfterHandler func(context.Context, *sql.DB) error

	AfterHandlers []AfterHandler
)
