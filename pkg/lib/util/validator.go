package util

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validate = sync.Pool{New: func() any {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("is_valid", func(fl validator.FieldLevel) bool {
		method := fl.Field().MethodByName("IsValid")
		if !method.IsValid() {
			return false
		}

		results := method.Call(nil)
		if len(results) == 0 {
			return false
		}

		if results[0].Kind() != reflect.Bool {
			return false
		}

		return results[0].Bool()
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

	for _, err := range errs {
		switch err.Tag() {

		case "is_valid":
			return fmt.Errorf("%s is invalid: %v", err.StructField(), err.Value())

		}
	}

	return nil
}

func Validate() *validator.Validate {
	return validate.Get().(*validator.Validate)
}
