package util

import "time"

var TimeEpochBegin = time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC)

func TimeSinceEpoch() time.Duration {
	return time.Now().UTC().Sub(TimeEpochBegin)
}
