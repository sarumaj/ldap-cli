package auth

import "testing"

func TestSchemeIsValid(t *testing.T) {
	for _, tt := range []struct {
		name string
		args Scheme
		want bool
	}{
		{"test#1", "ldap", true},
		{"test#2", "ldaps", true},
		{"test#3", "http", false},
		{"test#4", "ldapi", true},
		{"test#5", "cldap", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.IsValid()
			if got != tt.want {
				t.Errorf(`(%s).IsValid() failed: got: %t, want: %t`, tt.args, got, tt.want)
			}
		})
	}
}
