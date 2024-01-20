package util

import (
	"reflect"
	"testing"
)

func TestGetExecutablePath(t *testing.T) {
	got := GetExecutablePath()
	if got == "" {
		t.Errorf(`GetExecutablePath() failed: got : %q`, got)
	}
}

func TestRebuildStringSliceFlag(t *testing.T) {
	type args struct {
		flags     []string
		delimiter rune
	}

	type want struct {
		flags   []string
		wantErr bool
	}

	for _, tt := range []struct {
		name string
		args args
		want want
	}{
		{"test#1",
			args{[]string{`test#1;test#2`, `"test#3;test#4";test#5;`}, ';'},
			want{[]string{`test#1`, `test#2`, `test#3;test#4`, `test#5`}, false}},
		{"test#2", args{}, want{nil, true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RebuildStringSliceFlag(tt.args.flags, tt.args.delimiter)
			if (err != nil) != tt.want.wantErr {
				t.Errorf(`RebuildStringSliceFlag(..., %q) failed: %v`, tt.args.delimiter, err)
			}

			if err == nil && !reflect.DeepEqual(got, tt.want.flags) {
				t.Errorf(`RebuildStringSliceFlag(..., %q) failed: got: %v, want: %v`, tt.args.delimiter, got, tt.want.flags)
			}
		})
	}
}
