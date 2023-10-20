package auth

import (
	"strings"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

const (
	UNAUTHENTICATED = iota + 1
	SIMPLE
	MD5
	NTLM
)

var typeTranslation = map[AuthType]string{
	UNAUTHENTICATED: "UNAUTHENTICATED",
	SIMPLE:          "SIMPLE",
	MD5:             "MD5",
	NTLM:            "NTLM",
}

var _ util.ValidatorInterface = AuthType(0)

type AuthType int

// Validate type
func (t AuthType) IsValid() bool {
	switch t {

	case UNAUTHENTICATED, SIMPLE, MD5, NTLM:
		return true

	default:
		return false

	}
}

// Type as string
func (t AuthType) String() string {
	str, ok := typeTranslation[t]
	if ok {
		return str
	}

	return ""
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
