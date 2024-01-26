package auth

import (
	"fmt"
	"slices"
	"strings"

	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

const (
	UNAUTHENTICATED AuthType = iota + 1
	SIMPLE
	MD5
	NTLM
)

// typeTranslation is a map of AuthType to string
var typeTranslation = map[AuthType]string{
	UNAUTHENTICATED: "UNAUTHENTICATED",
	SIMPLE:          "SIMPLE",
	MD5:             "MD5",
	NTLM:            "NTLM",
}

var _ libutil.ValidatorInterface = AuthType(0)

// AuthType is an LDAP authentication type (RFC 4511)
type AuthType int

// IsValid returns true if the authentication type is valid
func (t AuthType) IsValid() bool {
	switch t {

	case UNAUTHENTICATED, SIMPLE, MD5, NTLM:
		return true

	default:
		return false

	}
}

// String returns the string representation of an authentication type
func (t AuthType) String() string {
	str, ok := typeTranslation[t]
	if ok {
		return str
	}

	return ""
}

func ListSupportedAuthTypes(quote bool) []string {
	var list []string
	for _, v := range typeTranslation {
		if quote {
			list = append(list, fmt.Sprintf("%q", v))
		} else {
			list = append(list, v)
		}
	}

	slices.Sort(list)
	return list
}

// Parse type from string
func TypeFromString(str string) AuthType {
	str = strings.ToUpper(str)
	for k, v := range typeTranslation {
		if strings.EqualFold(v, str) {
			return k
		}
	}

	return 0
}
