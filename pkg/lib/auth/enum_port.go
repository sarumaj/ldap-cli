package auth

type Port int

const (
	LDAP_RW  Port = 389
	LDAP_RO  Port = 3268
	LDAPS_RW Port = 636
	LDAPS_RO Port = 3269
)
