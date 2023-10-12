package auth

const (
	LDAP  Scheme = "ldap"
	LDAPS Scheme = "ldaps"
)

type Scheme string

func (s Scheme) IsValid() bool {
	switch s {

	case LDAP, LDAPS:
		return true

	default:
		return false

	}
}
