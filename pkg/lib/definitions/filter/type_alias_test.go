package filter

import (
	"reflect"
	"slices"
	"testing"
)

func TestListAliases(t *testing.T) {
	got, want := ListAliases(), make([]Alias, len(aliases))
	_ = copy(want, aliases)
	slices.SortStableFunc(want, func(a, b Alias) int {
		if a.Alias > b.Alias {
			return 1
		}

		return -1
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListAliases() failed: got: %v, want: %v", got, want)
	}
}
