package params

import (
	"strings"
)

// ConnParams specifies connection params.
// Source: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
var ConnParams = connParams{
	Database:               "database",
	Encrypt:                "encrypt",
	Password:               "password",
	ChangePassword:         "change password",
	UserID:                 "user id",
	Port:                   "port",
	TrustServerCertificate: "trustservercertificate",
	Certificate:            "certificate",
	TLSMin:                 "tlsmin",
	PacketSize:             "packet size",
	LogParam:               "log",
	ConnectionTimeout:      "connection timeout",
	HostNameInCertificate:  "hostnameincertificate",
	KeepAlive:              "keepalive",
	ServerSpn:              "serverspn",
	WorkstationID:          "workstation id",
	AppName:                "app name",
	ApplicationIntent:      "applicationintent",
	FailoverPartner:        "failoverpartner",
	FailOverPort:           "failoverport",
	DisableRetry:           "disableretry",
	Server:                 "server",
	Protocol:               "protocol",
	DialTimeout:            "dial timeout",
	Pipe:                   "pipe",
	MultiSubnetFailover:    "multisubnetfailover",
}

type connParams struct {
	Database               string
	Encrypt                string
	Password               string
	ChangePassword         string
	UserID                 string
	Port                   string
	TrustServerCertificate string
	Certificate            string
	TLSMin                 string
	PacketSize             string
	LogParam               string
	ConnectionTimeout      string
	HostNameInCertificate  string
	KeepAlive              string
	ServerSpn              string
	WorkstationID          string
	AppName                string
	ApplicationIntent      string
	FailoverPartner        string
	FailOverPort           string
	DisableRetry           string
	Server                 string
	Protocol               string
	DialTimeout            string
	Pipe                   string
	MultiSubnetFailover    string
}

func (p connParams) Exists(key string) bool {
	switch key := strings.ToLower(key); key {
	case p.Database,
		p.Encrypt,
		p.Password,
		p.ChangePassword,
		p.UserID,
		p.Port,
		p.TrustServerCertificate,
		p.Certificate,
		p.TLSMin,
		p.PacketSize,
		p.LogParam,
		p.ConnectionTimeout,
		p.HostNameInCertificate,
		p.KeepAlive,
		p.ServerSpn,
		p.WorkstationID,
		p.AppName,
		p.ApplicationIntent,
		p.FailoverPartner,
		p.FailOverPort,
		p.DisableRetry,
		p.Server,
		p.Protocol,
		p.DialTimeout,
		p.Pipe,
		p.MultiSubnetFailover:
		return true
	default:
		return false
	}
}
