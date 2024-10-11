package hypersql

import (
	"database/sql/driver"

	"github.com/qustavo/sqlhooks/v2"
)

type (
	DriverHooks []sqlhooks.Hooks

	DriverWrapper func(driver.Driver) driver.Driver

	DriverWrappers []DriverWrapper
)
