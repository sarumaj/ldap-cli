package auth

import (
	"strings"
	"testing"
)

func TestURLToBaseDirectoryPath(t *testing.T) {
	for _, tt := range []struct {
		name string
		args URL
		want string
	}{
		{"test#1",
			URL{Host: "example.com"},
			"DC=example,DC=com"},
		{"test#2",
			URL{Host: ".example.com"},
			"DC=example,DC=com"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.ToBaseDirectoryPath()
			if got != tt.want {
				t.Errorf(`(%v).Validate() failed: got: %q, want: %q`, tt.args, got, tt.want)
			}
		})
	}
}

func TestURLBuild(t *testing.T) {
	for _, tt := range []struct {
		name string
		args *URL
		want string
	}{
		{"test#1",
			NewURL().SetScheme(LDAPS).SetHost("example.com").SetPort(LDAPS_RW),
			"ldaps://example.com:636"},
	} {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.args.String(); got != tt.want {
				t.Errorf(`(%v).String() failed: got: %q, want: %q`, tt.args, got, tt.want)
			}

			if got := tt.args.HostPort(); got != strings.TrimPrefix(tt.want, string(tt.args.Scheme)+"://") {
				t.Errorf(`(%v).HostPort() failed: got: %q, want: %q`, tt.args, got, tt.want)
			}
		})
	}
}

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
			URL{Scheme: LDAPS},
			true},
		{"test#3",
			URL{Scheme: LDAPS, Host: "example.com"},
			true},
		{"test#4",
			URL{Scheme: LDAPS, Host: "example.com", Port: LDAPS_RW},
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
