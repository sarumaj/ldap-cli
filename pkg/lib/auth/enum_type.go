package auth

import (
	"strings"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

const (
	SIMPLE = iota + 1
	MD5
	NTLM
)

var typeTranslation = map[Type]string{
	SIMPLE: "SIMPLE",
	MD5:    "MD5",
	NTLM:   "NTLM",
}

var _ util.ValidatorInterface = Type(0)

type Type int

// Validate type
func (t Type) IsValid() bool {
	switch t {

	case SIMPLE, MD5, NTLM:
		return true

	default:
		return false

	}
}

// Type as string
func (t Type) String() string {
	str, ok := typeTranslation[t]
	if ok {
		return str
	}

	return ""
}

// Parse type from string
func TypeFromString(str string) Type {
	str = strings.ToUpper(str)
	for k, v := range typeTranslation {
		if strings.EqualFold(v, str) {
			return k
		}
	}

	return 0
}
