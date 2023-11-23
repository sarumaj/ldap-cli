package attributes

import (
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	for _, tt := range []struct {
		name string
		args Map
		want Attributes
	}{
		{"test#1", Map{name: "", accountExpires: 0}, Attributes{accountExpires, name}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.Keys()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(Map).Keys() failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}

}
