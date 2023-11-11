package auth

import (
	"reflect"
	"strings"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type Type int

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

func (t Type) IsValid() bool {
	switch t {

	case SIMPLE, MD5, NTLM:
		return true

	default:
		return false

	}
}

func (t Type) String() string {
	str, ok := typeTranslation[t]
	if ok {
		return str
	}

	return ""
}

func init() {
	util.Validate().RegisterCustomTypeFunc(func(field reflect.Value) any {
		v, ok := field.Interface().(Type)
		if ok {
			return v.String()
		}

		return nil
	}, Type(0))
}

func TypeFromString(str string) Type {
	str = strings.ToUpper(str)
	for k, v := range typeTranslation {
		if v == str {
			return k
		}
	}

	return 0
}
