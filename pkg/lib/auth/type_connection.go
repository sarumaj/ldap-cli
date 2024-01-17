package auth

import ldap "github.com/go-ldap/ldap/v3"

// Connection object
type Connection struct {
	// LDAP connection
	*ldap.Conn
	// Dial options
	*DialOptions
	// Remote host
	remoteHost string
}

// Close closes the underlying TCP connection
func (c Connection) Close() error { return c.Conn.Close() }

// RemoteHost returns the remote host of the domain controller
func (c Connection) RemoteHost() string { return c.remoteHost }
