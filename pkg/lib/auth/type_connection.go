package auth

import ldap "github.com/go-ldap/ldap/v3"

// Connection object
type Connection struct {
	*ldap.Conn
	*DialOptions
	remoteHost string
}

func (c Connection) Close() error       { return c.Conn.Close() } // Close TCP connection
func (c Connection) RemoteHost() string { return c.remoteHost }   // retrieve address of current domain controller
