package attributes

import (
	"reflect"
	"testing"
)

func TestAttributesToStringSlice(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []Attribute
		want []string
	}{
		{"test#1", nil, nil},
		{"test#2",
			[]Attribute{AttributeAccountExpires, AttributeEnabled, AttributeCommonName},
			[]string{"accountExpires", "cn", "userAccountControl"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := AttributesToStringSlice(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`AttributesToStringSlice(...) failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}
