package auth

import "testing"

func TestTypeIsValid(t *testing.T) {
	for _, tt := range []struct {
		name string
		args AuthType
		want bool
	}{
		{"test#1", SIMPLE, true},
		{"test#2", NTLM, true},
		{"test#3", 0, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.IsValid()
			if got != tt.want {
				t.Errorf(`(%q).IsValid() failed: got: %t, want: %t`, tt.args, got, tt.want)
			}
		})
	}

}

func TestTypeFromString(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want AuthType
	}{
		{"test#1", "simple", SIMPLE},
		{"test#2", "ntlm", NTLM},
		{"test#3", "invalid", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := TypeFromString(tt.args)
			if got != tt.want {
				t.Errorf(`TypeFromString(%q) failed: got: %q, want: %q`, tt.args, got, tt.want)
			}
		})
	}
}
