package filter

import (
	"reflect"
	"testing"
)

func Test_balanceClosingParenthesis(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want string
	}{
		{"test#1", `(test))((test)))`, `(test))((test)))`},
		{"test#2", `(test))))test)))`, `(test))))test)))`},
		{"test#3", `((test))((test))`, `((test))((test))`},
		{"test#4", `test))((test)))`, `(test))((test))))`},
		{"test#5", `((test))((test`, `((test))((test`},
		{"test#6", `test))((((test`, `(test))((((test))))`},
		{"test#7", `test))test`, `(test))test)`},
		{"test#8", `test`, `(test)`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := balanceClosingParenthesis(tt.args); got != tt.want {
				t.Errorf(`Pad(%q) failed: got: %q, want: %q`, tt.args, got, tt.want)
			}
		})
	}
}

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
		{"test#2", `(&(|(|(CN=test#1)(CN=test#2))(CN=test#1))(CN=test#2))`, want{`(&(|(|(CN=test#1)(CN=test#2))(CN=test#1))(CN=test#2))`, false}},
		{"test#3", `(&(CN=test#1)(CN=test#2))`, want{`(&(CN=test#1)(CN=test#2))`, false}},
		{"test#4", `(|(!(CN=test#1))(Unknown=test#2))`, want{`(|(!(CN=test#1))(Unknown=test#2))`, false}},
		{"test#5", `(&(!(CN=test#1))(CN=test#2)`, want{"", true}},
		{"test#6", `(&(!(CN=test#1))(CN:test#2))`, want{"", true}},
		{"test#7", `(CN:test)`, want{"", true}},
		{"test#8", `(&(ObjectClass=user)(!(UserAccountControl:$BAND:=2)))`,
			want{(`(&` +
				`(ObjectClass=user)` +
				(`(!` +
					`(UserAccountControl:1.2.840.113556.1.4.803:=2)`) +
				`)`) +
				`)`, false}},
		{"test#9", `(&$ID(1)$ID(2))`,
			want{(`(&` +
				(`(|` +
					`(CN=1)` +
					`(DisplayName=1)` +
					(`(|` +
						`(DistinguishedName=1)` +
						`(DN=1)`) +
					`)` +
					`(Name=1)` +
					`(SAMAccountName=1)` +
					`(UserPrincipalName=1)`) +
				`)` +
				(`(|` +
					`(CN=2)` +
					`(DisplayName=2)` +
					(`(|` +
						`(DistinguishedName=2)` +
						`(DN=2)`) +
					`)` +
					`(Name=2)` +
					`(SAMAccountName=2)` +
					`(UserPrincipalName=2)`) +
				`)`) +
				`)`, false}},
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

func Test_splitFilterComponents(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want []string
	}{
		{"test#1", `(test))((test)))`, []string{"(test)", ")(", "(test)"}},
		{"test#2", `(test))))test)))`, []string{"(test)"}},
		{"test#3", `((test))((test))`, []string{"((test))", "((test))"}},
		{"test#4", `((test))((test)))`, []string{"((test))", "((test))"}},
		{"test#5", `((test))((test)))))`, []string{"((test))", "((test))"}},
		{"test#6", `((test))((((test))`, []string{"((test))"}},
		{"test#7", `test))test`, nil},
		{"test#8", `test((test`, nil},
		{"test#9", `(test)`, []string{"(test)"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitFilterComponents(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`split(%q) failed: got: %q, want: %q`, tt.args, got, tt.want)
			}
		})
	}
}
