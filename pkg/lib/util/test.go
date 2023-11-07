package util

import (
	"os"
	"testing"
)

func SkipOAT(t testing.TB) {
	switch os.Getenv("TEST_OAT") {

	case "true", "True", "TRUE", "1", "y", "yes", "YES":
		return

	}

	t.Skipf("Running only FAT tests, skipping %q, since it requires extensive mock-up", t.Name())
}
