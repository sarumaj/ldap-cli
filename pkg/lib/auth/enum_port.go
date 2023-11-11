package auth

type Port int

const (
	// Local catalogue port
	LDAP_RW Port = 389
	// Global catalogue port
	LDAP_RO Port = 3268
	// Local catalogue port over TLS
	LDAPS_RW Port = 636
	// Global catalogue port over TLS
	LDAPS_RO Port = 3269
)
