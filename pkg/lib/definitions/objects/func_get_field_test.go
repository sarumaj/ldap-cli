package objects

import (
	"reflect"
	"testing"
)

func TestGetField(t *testing.T) {
	type object struct {
		String string
		Int    int
		Int64  int64
	}
	type args struct {
		o object
		p string
	}
	for _, tt := range []struct {
		name string
		args args
		want any
	}{
		{"test#1", args{object{"test", 1, 2}, "String"}, "test"},
		{"test#2", args{object{"test", 1, 2}, "Int"}, 1},
		{"test#3", args{object{"test", 1, 2}, "Int64"}, int64(2)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := GetField[any](&tt.args.o, tt.args.p)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`GetField[any](..., %q) failed: got: %v, want: %v`, tt.args.p, got, tt.want)
			}
		})
	}
}
