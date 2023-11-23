package util

import (
	"time"

	"bou.ke/monkey"
)

func PatchForTimeNow() *monkey.PatchGuard {
	return monkey.Patch(
		time.Now,
		func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) },
	)
}
