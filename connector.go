package hypersql

import (
	"context"
	"database/sql/driver"
)

var _ driver.Connector = (*dsnConnector)(nil)

type dsnConnector struct {
	dsn    string
	driver driver.Driver
}

func (c *dsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return c.driver.Open(c.dsn)
}

func (c *dsnConnector) Driver() driver.Driver {
	return c.driver
}

type wrapConnector struct {
	c driver.Connector
}

func (w *wrapConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return w.c.Connect(ctx)
}

func (w *wrapConnector) Driver() driver.Driver {
	return w.c.Driver()
}

var _ driver.Connector = (*wrapConnector)(nil)

func WrapConnector(c driver.Connector) driver.Connector {
	return &wrapConnector{c: c}
}
