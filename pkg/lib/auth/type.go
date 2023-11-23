package auth

import (
	"strings"
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

func TypeFromString(str string) Type {
	str = strings.ToUpper(str)
	for k, v := range typeTranslation {
		if v == str {
			return k
		}
	}

	return 0
}
