package util

import (
	"regexp"
	"testing"
)

func TestNewGUID(t *testing.T) {
	got := NewGUID()
	if guidRegex := regexp.MustCompile(`^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$`); !guidRegex.MatchString(got) {
		t.Errorf(`NewGUID() failed: %q does not match %q`, got, guidRegex)
	}
}
