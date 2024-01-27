package util

import (
	"os"
	"testing"
)

func TestSkipOAT(t *testing.T) {
	t.Cleanup(func() { _ = os.Unsetenv("TEST_OAT") })

	_ = os.Setenv("TEST_OAT", "true")
	SkipOAT(t)

	_ = os.Setenv("TEST_OAT", "false")
	SkipOAT(t)
}
