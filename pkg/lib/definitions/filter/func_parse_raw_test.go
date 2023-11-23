package filter

import "testing"

func TestParseRaw(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    string
		wantErr bool
	}{
		{"test#1", `(cn=test)`, false},
		{"test#2", `(&(!(cn=test#1))(cn=test#2))`, false},
		{"test#3", `(|(!(cn=test#1))(unknown=test#2))`, false},
		{"test#4", `(&(!(cn=test#1))(cn=test#2)`, true},
		{"test#5", `(&(!(cn=test#1))(cn:test#2))`, true},
		{"test#6", `(cn:test)`, true},
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
