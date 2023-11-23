package auth

import "github.com/go-ldap/ldap/v3"

// Connection object
type Connection struct {
	connection *ldap.Conn
	options    *DialOptions
	remoteHost string
}

func (c Connection) Close() error { return c.connection.Close() }
