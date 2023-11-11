package util

import "testing"

func TestToTitleNoLower(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want string
	}{
		{"test#1", "HELLO WORLD", "HELLO WORLD"},
		{"test#2", "hello world", "Hello World"},
		{"test#3", "helloWorld", "HelloWorld"},
	} {
		got := ToTitleNoLower(tt.args)
		if got != tt.want {
			t.Errorf(`ToTitleNoLower(%q) failed: got: %q, want: %q`, tt.args, got, tt.want)
		}
	}
}
