package auth

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

// For internal usage
var validate = func() *validator.Validate {
	validate := util.Validate()

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
var _ util.ValidatorInterface = typeString("")

type typeString string

func (t typeString) IsValid() bool { return TypeFromString(string(t)).IsValid() }
