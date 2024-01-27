package util

import (
	"math/big"
	"time"
)

// Time1601 returns the epoch start date: 1601-01-01 00:00:00 UTC
func Time1601() time.Time { return time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC) }

// TimeAfter1601 returns time after 1601-01-01 00:00:00 UTC. Offset should be in 0.1 µs
func TimeAfter1601(offset int64) time.Time {
	// µs since UNIX epoch
	begin := big.NewInt(Time1601().UnixMicro())

	// offset in 0.1 µs
	elapsed := big.NewInt(offset)

	// convert to µs
	µs, rem := elapsed.DivMod(elapsed, big.NewInt(10), big.NewInt(10))

	// add µs, get elapsed ns
	sum, ns := begin.Add(begin, µs).Int64(), rem.Mul(rem, big.NewInt(1000)).Int64()

	return time.UnixMicro(sum).Add(time.Duration(ns)).UTC()
}

// TimeSince1601 returns time since 1601-01-01 00:00:00 UTC
func TimeSince1601() time.Duration { return Now().UTC().Sub(Time1601()) }
