package attributes

import (
	"reflect"
	"testing"
)

func TestAttributesToAttributeList(t *testing.T) {
	for _, tt := range []struct {
		name string
		args Attributes
		want []string
	}{
		{"test#1", nil, nil},
		{"test#2",
			[]Attribute{AttributeAccountExpires(), AttributeUserAccountControl(), AttributeCommonName()},
			[]string{"accountexpires", "cn", "useraccountcontrol"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.ToAttributeList()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`AttributesToStringSlice(...) failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}
