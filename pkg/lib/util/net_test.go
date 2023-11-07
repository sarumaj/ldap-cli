package util

import (
	"testing"
)

func TestLookupAddress(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want string
	}{
		{"test#1", "127.0.0.1:443", "localhost:443"},
		{"test#2", "127.0.0.1", "localhost"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := LookupAddress(tt.args)
			if got != tt.want {
				t.Errorf(`LookupAddress(%q) failed: got: %q, want: %q`, tt.args, got, tt.want)
			}
		})
	}
}
