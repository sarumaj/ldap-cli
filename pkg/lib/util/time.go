package util

import (
	"math/big"
	"time"
)

func Time1601() time.Time { return time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC) }

// Retrieve date by applying offset. It should be in 0.1 µs
func TimeAfter1601(offset int64) time.Time {
	// µs since UNIX epoch
	begin := big.NewInt(Time1601().UnixMicro())

	// offset in 0.1 µs
	elapsed := big.NewInt(offset)

	// convert to µs
	µs, ns := elapsed.DivMod(elapsed, big.NewInt(10), big.NewInt(10))

	if r := begin.Add(begin, µs); r.IsInt64() {
		return time.UnixMicro(r.Int64()).Add(time.Duration(ns.Mul(ns, big.NewInt(1000)).Int64()))
	}

	// handle overflow
	return time.UnixMicro(1<<63 - 1)
}

func TimeSince1601() time.Duration {
	return time.Now().UTC().Sub(Time1601())
}
