package util

import (
	"os"
	"testing"
	"time"
)

// Exit is reference to os.Exit (can be mocked)
var Exit = os.Exit

// Now is reference to time.Now (can be mocked)
var Now = time.Now

// SkipOAT skips the test if TEST_OAT is not set to true
func SkipOAT(t testing.TB) {
	switch os.Getenv("TEST_OAT") {

	case "true", "True", "TRUE", "1", "y", "yes", "YES":
		return

	}

	t.Skipf("Running only FAT tests, skipping %q, since it requires extensive mock-up", t.Name())
}
