package util

import "testing"

func TestHexifyUnhexify(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want string
	}{
		{"test#1", "test", "\\x74\\x65\\x73\\x74"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Hexify(tt.args)
			if encoded != tt.want {
				t.Errorf(`hexify(%q) failed: got: %q, want: %q`, tt.args, encoded, tt.want)
			}

			decoded := Unhexify(encoded)
			if decoded != tt.args {
				t.Errorf(`Unhexify(%q) failed: got: %q, want: %q`, encoded, decoded, tt.args)
			}
		})
	}

}
