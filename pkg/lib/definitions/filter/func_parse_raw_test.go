package filter

import "testing"

func TestParseRaw(t *testing.T) {
	type want struct {
		result string
		err    bool
	}
	for _, tt := range []struct {
		name string
		args string
		want want
	}{
		{"test#1", `(CN=test)`, want{`(CN=test)`, false}},
		{"test#2", `(&(!(CN=test#1))(CN=test#2))`, want{`(&(!(CN=test#1))(CN=test#2))`, false}},
		{"test#3", `(|(!(CN=test#1))(Unknown=test#2))`, want{`(|(!(CN=test#1))(Unknown=test#2))`, false}},
		{"test#4", `(&(!(CN=test#1))(CN=test#2)`, want{"", true}},
		{"test#5", `(&(!(CN=test#1))(CN:test#2))`, want{"", true}},
		{"test#6", `(CN:test)`, want{"", true}},
		{"test#7", `(&(ObjectClass=user)(!(UserAccountControl:${AND}:=2)))`, want{`(&(ObjectClass=user)(!(UserAccountControl:1.2.840.113556.1.4.803:=2)))`, false}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRaw(tt.args)
			if (err == nil) == tt.want.err {
				t.Errorf(`ParseRaw(%q) failed: %v`, tt.args, err)
			}

			if got != nil && got.String() != tt.want.result {
				t.Errorf(`ParseRaw(%q) failed: got: %q, want: %q`, tt.args, got, tt.want.result)
			}
		})
	}
}
