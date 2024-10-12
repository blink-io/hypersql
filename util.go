package hypersql

import (
	"context"
	"database/sql/driver"
	"errors"
	"net"
	"strconv"
)

// DoPingContext does invoke ping(context.Context).
// Ignore driver.ErrSkip when the Conn does not implement driver.Pinger interface.
func DoPingContext(ctx context.Context, pinger Pinger) error {
	if err := pinger.PingContext(ctx); err != nil && !errors.Is(err, driver.ErrSkip) {
		return err
	}
	return nil
}

func hostPortToAddr(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
