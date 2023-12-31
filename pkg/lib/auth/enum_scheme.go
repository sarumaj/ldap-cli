package auth

import libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"

const (
	LDAP  Scheme = "ldap"
	LDAPS Scheme = "ldaps"
)

var _ libutil.ValidatorInterface = Scheme("")

type Scheme string

// Validate scheme
func (s Scheme) IsValid() bool {
	switch s {

	case LDAP, LDAPS:
		return true

	default:
		return false

	}
}
