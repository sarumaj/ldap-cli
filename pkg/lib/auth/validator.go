package auth

import (
	"reflect"

	validator "github.com/go-playground/validator/v10"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

// For internal usage
var validate = func() *validator.Validate {
	validate := libutil.Validate()

	// custom type validator for Type
	// string is casted to typeString which implements util.ValidatorInterface
	validate.RegisterCustomTypeFunc(
		func(field reflect.Value) any { return typeString(field.Interface().(AuthType).String()) },
		AuthType(0),
	)

	return validate
}()

/*
 * String as util.ValidatorInterface for Type
 */

var _ libutil.ValidatorInterface = typeString("")

// typeString is a string that implements util.ValidatorInterface
type typeString string

// IsValid returns true if the string is a valid Type
func (t typeString) IsValid() bool { return TypeFromString(string(t)).IsValid() }
