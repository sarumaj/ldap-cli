package auth

import (
	"testing"
)

func TestURLValidate(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    URL
		wantErr bool
	}{
		{"test#1",
			*NewURL(),
			true},
		{"test#2",
			*(NewURL().SetScheme(LDAPS)),
			true},
		{"test#3",
			*(NewURL().SetScheme(LDAPS).SetHost("example.com")),
			true},
		{"test#4",
			*(NewURL().SetScheme(LDAPS).SetHost("example.com").SetPort(LDAPS_RW)),
			false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := &tt.args
			err := opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf(`(%v).Validate() failed: %v`, tt.args, err)
			}
		})
	}
}

func TestURLFromString(t *testing.T) {
	type want struct {
		URL     *URL
		wantErr bool
	}
	for _, tt := range []struct {
		name string
		args string
		want want
	}{
		{"test#1", "https://example.com", want{nil, true}},
		{"test#2", "https://example.com:8080", want{nil, true}},
		{"test#3", "ldaps://example.com:8080", want{&URL{LDAPS, "example.com", 8080}, false}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLFromString(tt.args)
			if (err != nil) != tt.want.wantErr {
				t.Errorf(`URLFromString(%q) failed: %v`, tt.args, err)
			}
			if got != nil && got.String() != tt.want.URL.String() {
				t.Errorf(`URLFromString(%q) failed: got: %q, want: %q`, tt.args, got, tt.want.URL)
			}
		})
	}
}
