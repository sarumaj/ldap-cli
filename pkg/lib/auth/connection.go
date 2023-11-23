package auth

import "github.com/go-ldap/ldap/v3"

type Connection struct {
	connection   *ldap.Conn
	options      *DialOptions
	remoteHost   string
	retryCounter int
}

func (c Connection) Close() error { return c.connection.Close() }
