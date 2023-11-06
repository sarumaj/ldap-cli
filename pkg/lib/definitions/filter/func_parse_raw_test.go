package filter

import "testing"

func TestParseRaw(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    string
		wantErr bool
	}{
		{"test#1", `(CN=test)`, false},
		{"test#2", `(&(!(CN=test#1))(CN=test#2))`, false},
		{"test#3", `(|(!(CN=test#1))(Unknown=test#2))`, false},
		{"test#4", `(&(!(CN=test#1))(CN=test#2)`, true},
		{"test#5", `(&(!(CN=test#1))(CN:test#2))`, true},
		{"test#6", `(CN:test)`, true},
		{"test#7", `(&(ObjectClass=user)(!(UserAccountControl:1.2.840.113556.1.4.803:=2)))`, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRaw(tt.args)
			if (err == nil) == tt.wantErr {
				t.Errorf(`ParseRaw(%q) failed: %v`, tt.args, err)
			}

			if got != nil && got.String() != tt.args {
				t.Errorf(`ParseRaw(%[2]q) failed: got: %[1]q, want: %[2]q`, got, tt.args)
			}
		})
	}
}
