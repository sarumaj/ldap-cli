package util

import (
	"math/big"
	"time"
)

func Time1601() time.Time { return time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC) }

// Retrieve date by applying offset. It should be in 0.1 µs
func TimeAfter1601(offset int64) time.Time {
	begin := Time1601()

	// µs since UNIX epoch
	n := big.NewInt(begin.UnixMicro())

	// offset in 0.1 µs
	n2 := big.NewInt(offset)

	// convert to µs
	n2 = n2.Div(n2, big.NewInt(10))

	// result
	r := n.Add(n, n2)

	if r.IsInt64() {
		return time.UnixMicro(r.Int64())
	}

	// handle overflow
	return time.UnixMicro(1<<63 - 1)
}

func TimeSince1601() time.Duration {
	return time.Now().UTC().Sub(Time1601())
}
