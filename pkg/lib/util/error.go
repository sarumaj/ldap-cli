package util

import "errors"

var ErrBindFirst = errors.New("unbound request, bind to LDAP server first")

func PanicIfError[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
