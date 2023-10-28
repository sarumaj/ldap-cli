package attributes

import (
	"reflect"
	"testing"
)

func TestGroupTypeEval(t *testing.T) {
	for _, tt := range []struct {
		name string
		args FlagsetGroupType
		want []string
	}{
		{"test#1",
			GROUP_TYPE_APP_BASIC | GROUP_TYPE_SECURITY | GROUP_TYPE_UNIVERSAL,
			[]string{"APP_BASIC", "SECURITY", "UNIVERSAL"}},
		{"test#2",
			GROUP_TYPE_APP_BASIC | GROUP_TYPE_UNIVERSAL,
			[]string{"APP_BASIC", "DISTRIBUTION", "UNIVERSAL"}},
		{"test#3", 0, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.Eval()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(GroupType).Eval() failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}
