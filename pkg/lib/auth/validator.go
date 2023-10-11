package auth

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

var validate = func() *validator.Validate {
	validate := util.Validate()

	validate.RegisterCustomTypeFunc(func(field reflect.Value) any {
		v, ok := field.Interface().(Type)
		if ok {
			return v.String()
		}

		return nil
	}, Type(0))

	return validate
}()
