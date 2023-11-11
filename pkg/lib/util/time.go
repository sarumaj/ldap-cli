package util

import (
	"math/big"
	"time"
)

func Time1601() time.Time { return time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC) }

// Retrieve date by applying offset. It should be in 0.01 ns
func TimeAfter1601(offset int64) time.Time {
	begin := Time1601()

	// nanoseconds since UNIX epoch
	n := big.NewInt(begin.UnixNano())

	// offset in ns
	n2 := big.NewInt(offset)
	n2 = n2.Mul(n2, big.NewInt(100))

	// result
	r := n.Add(n, n2)

	if r.IsInt64() {
		return time.Unix(0, r.Int64())
	}

	// handle overflow
	return time.Unix(0, 1<<63-1)
}

func TimeSince1601() time.Duration {
	return time.Now().UTC().Sub(Time1601())
}
