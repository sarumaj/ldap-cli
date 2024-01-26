package auth

import libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"

const (
	LDAP  Scheme = "ldap"
	LDAPS Scheme = "ldaps"
)

var _ libutil.ValidatorInterface = Scheme("")

// Scheme is an LDAP scheme
type Scheme string

// IsValid returns true if the scheme is valid
func (s Scheme) IsValid() bool {
	switch s {

	case LDAP, LDAPS:
		return true

	default:
		return false

	}
}
