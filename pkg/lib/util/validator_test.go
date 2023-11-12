package util

import (
	"strings"
	"testing"
)

type invalid struct {
	SomeString           string        `validate:"required"`
	SomeValidator        implValidator `validate:"is_valid"`
	SomeWannaBeValidator string        `validate:"is_valid"`
	SomeInt              int           `validate:"gt=0"`
	SomeFloat            float64       `validate:"required_if=SomeInt 0"`
	SomeBool             bool          `validate:"required_unless=SomeInt 1"`
	SomePassing          bool          `validate:"required"`
	SomeTiger            string        `validate:"tiger192"`
}

type implValidator string

func (implValidator) IsValid() bool { return false }

func TestFormatError(t *testing.T) {
	validator := Validate()
	err := FormatError(validator.Struct(&invalid{SomeValidator: "empty", SomePassing: true}))
	t.Log(err.Error())
	for _, msg := range []string{
		`"SomeString" is required`,
		`"SomeValidator" is invalid: empty`,
		`"SomeInt" should be greater than 0`,
		`"SomeFloat" is required, if "SomeInt" is "0"`,
		`"SomeBool" is required, unless "SomeInt" is "1"`,
		`Key: 'invalid.SomeTiger' Error:Field validation for 'SomeTiger' failed on the 'tiger192' tag`,
	} {
		if !strings.Contains(err.Error(), msg) {
			t.Errorf(`FormatError(validator.Struct(&invalid{})) failed: want: %q`, msg)
		}
	}
}
