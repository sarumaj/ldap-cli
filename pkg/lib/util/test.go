package util

import (
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
)

func PatchForTimeNow() *monkey.PatchGuard {
	return monkey.Patch(
		time.Now,
		func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) },
	)
}

func SkipOAT(t testing.TB) {
	switch os.Getenv("TEST_OAT") {

	case "true", "True", "TRUE", "1", "y", "yes", "YES":
		return

	}

	t.Skipf("Running only FAT tests, skipping %q, since it requires extensive mock-up", t.Name())
}
