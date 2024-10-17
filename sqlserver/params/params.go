package params

import (
	"strings"
)

// ConnParams specifies connection params.
// Source: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
var ConnParams = connParams{
	Certificate:       "certificate",
	KeepAlive:         "keepAlive",
	ServerSPN:         "ServerSPN",
	ApplicationIntent: "ApplicationIntent",
	Workstation:       "Workstation",
	Instance:          "instance",
}

type connParams struct {
	Certificate       string
	KeepAlive         string
	ServerSPN         string
	ApplicationIntent string
	Workstation       string
	Instance          string
}

func (p connParams) Exists(key string) bool {
	switch key := strings.ToLower(key); key {
	case p.KeepAlive,
		p.Certificate,
		p.Workstation,
		p.ServerSPN:
		return true
	default:
		return false
	}
}
