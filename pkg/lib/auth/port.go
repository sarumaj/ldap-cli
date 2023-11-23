package auth

type Port uint

const (
	LDAP         Port = 389
	LDAP_GLOBAL  Port = 3268
	LDAPS        Port = 636
	LDAPS_GLOBAL Port = 3269
)
