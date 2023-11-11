package attributes

import (
	"reflect"
	"testing"
)

func TestSAMAccountTypeEval(t *testing.T) {
	for _, tt := range []struct {
		name string
		args SAMAccountType
		want []string
	}{
		{"test#1",
			SAM_ACCOUNT_TYPE_USER_OBJECT,
			[]string{"NORMAL_USER_ACCOUNT", "USER_OBJECT"}},
		{"test#2",
			SAM_ACCOUNT_TYPE_DOMAIN_OBJECT,
			[]string{"DOMAIN_OBJECT"}},
		{"test#3", 0, []string{"DOMAIN_OBJECT"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.Eval()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(SAMAccountType).Eval() failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}
