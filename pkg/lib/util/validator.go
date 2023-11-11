package util

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type ValidatorInterface interface{ IsValid() bool }

var validate = sync.Pool{New: func() any {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("is_valid", func(fl validator.FieldLevel) bool {
		v, ok := fl.Field().Interface().(ValidatorInterface)
		if ok {
			return v.IsValid()
		}

		return false
	})

	return validate
}}

func FormatError(err error) error {
	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var newErrs []error
	for _, err := range errs {
		switch err.Tag() {

		case "gt", "gte", "lt", "lte":
			condition := map[string]string{
				"gt":  "greater than",
				"gte": "greater than or equal to",
				"lt":  "less than",
				"lte": "less than or equal to",
			}[err.Tag()]
			newErrs = append(newErrs, fmt.Errorf("%q should be %s %s", err.StructField(), condition, err.Param()))

		case "is_valid":
			newErrs = append(newErrs, fmt.Errorf("%q is invalid: %v", err.StructField(), err.Value()))

		case "required":
			newErrs = append(newErrs, fmt.Errorf("%q is required", err.StructField()))

		case "required_if", "required_unless":
			condition := strings.TrimPrefix(err.Tag(), "required_")
			field, value, found := strings.Cut(err.Param(), " ")
			if found {
				newErrs = append(newErrs, fmt.Errorf("%q is required, %s %q is %q", err.StructField(), condition, field, value))
			} else {
				newErrs = append(newErrs, fmt.Errorf("%q is required, %s %q", err.StructField(), condition, err.Param()))
			}

		default:
			newErrs = append(newErrs, err)

		}
	}

	return errors.Join(newErrs...)
}

func Validate() *validator.Validate {
	return validate.Get().(*validator.Validate)
}
