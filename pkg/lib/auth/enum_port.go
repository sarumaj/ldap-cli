package auth

type Port int

const (
	// LDAP_RW is the default LDAP port (local catalogue port)
	LDAP_RW Port = 389
	// LDAP_RO is the global catalogue port
	LDAP_RO Port = 3268
	// LDAPS_RW is the default LDAP port over TLS (local catalogue port)
	LDAPS_RW Port = 636
	// LDAPS_RO is the global catalogue port over TLS
	LDAPS_RO Port = 3269
)
