package util

import (
	"testing"
	"time"

	"bou.ke/monkey"
)

func TestTimeAfter1601(t *testing.T) {
	for _, tt := range []struct {
		name string
		args int64
		want time.Time
	}{
		{"test#1", 128271382742968750, time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)},
		{"test#2", 1<<63 - 1, time.Date(30828, 9, 14, 2, 48, 5, 477587000, time.UTC)},
	} {
		t.Log(Time1601.Add(time.Duration(tt.args*100) * time.Nanosecond))
		t.Run(tt.name, func(t *testing.T) {
			got := TimeAfter1601(tt.args)
			if !got.Equal(tt.want) {
				t.Errorf(`TimeAfter1601(%d) failed: got: %s, want: %s`, tt.args, got, tt.want)
			}
		})
	}
}

func TestTimeSince1601(t *testing.T) {
	defer monkey.Patch(
		time.Now,
		func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) },
	).Unpatch()

	for _, tt := range []struct {
		name string
		want time.Duration
	}{
		{"test#1", func() time.Duration { d, _ := time.ParseDuration("2562047h47m16.854775807s"); return d }()},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeSince1601()
			if got != tt.want {
				t.Errorf(`TimeSince1601() failed: got: %s, want: %s`, got, tt.want)
			}
		})
	}
}
