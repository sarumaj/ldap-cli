package attributes

import (
	"reflect"
	"testing"
)

func TestUserAccountControlEval(t *testing.T) {
	for _, tt := range []struct {
		name string
		args FlagsetUserAccountControl
		want []string
	}{
		{"test#1",
			USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE | USER_ACCOUNT_CONTROL_NORMAL_ACCOUNT | USER_ACCOUNT_CONTROL_PASSWORD_CANT_CHANGE,
			[]string{"ACCOUNT_DISABLE", "NORMAL_ACCOUNT", "PASSWORD_CANT_CHANGE"}},
		{"test#2", 0, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.Eval()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(UserAccountControl).Eval() failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}

func TestUserAccountControlString(t *testing.T) {
	for _, tt := range []struct {
		name string
		args FlagsetUserAccountControl
		want []string
	}{
		{"test#1",
			USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE | USER_ACCOUNT_CONTROL_NORMAL_ACCOUNT | USER_ACCOUNT_CONTROL_PASSWORD_CANT_CHANGE,
			[]string{"ACCOUNT_DISABLE", "NORMAL_ACCOUNT", "PASSWORD_CANT_CHANGE"}},
		{"test#2", 0, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.Eval()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(UserAccountControl).Eval() failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}
