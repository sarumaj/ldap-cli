package objects

import (
	"reflect"
	"testing"

	diff "github.com/r3labs/diff/v3"
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

func TestHexifyUnhexify(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
	}{
		{"test#1", "test"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			encoded := hexify(tt.args)
			decoded := Unhexify(encoded)
			if decoded != tt.args {
				t.Errorf(`hexify(%[1]q) failed: got: %[2]q, want: %[1]q`, tt.args, decoded)
			}
		})
	}

}

func TestReadMap(t *testing.T) {
	type object struct {
		Strings []string
		Bool    bool
		Int     int
		Int64   int64
		String  string
	}
	type args struct {
		o object
		m map[string]any
	}
	for _, tt := range []struct {
		name string
		args args
		want object
	}{
		{"test#1", args{
			object{},
			map[string]any{
				"strings": "test",
				"bool":    "true",
				"int":     "12",
				"int64":   "13",
				"string":  []string{"test#1.1", "test#1.2"},
			},
		}, object{
			Strings: []string{"test"},
			Bool:    true,
			Int:     12,
			Int64:   13,
			String:  "test#1.1;test#1.2",
		}},
		{"test#2", args{
			object{},
			map[string]any{
				"strings": []string{"test#1.1", "test#1.2"},
				"bool":    "true",
				"int":     "12",
				"int64":   "13",
				"string":  "test",
			},
		}, object{
			Strings: []string{"test#1.1", "test#1.2"},
			Bool:    true,
			Int:     12,
			Int64:   13,
			String:  "test",
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.o
			readMap(&got, tt.args.m)

			changelogs, err := diff.Diff(got, tt.want, diff.DisableStructValues())
			if err != nil {
				t.Error(err)
			}

			for _, changelog := range changelogs {
				t.Errorf(`readMap(...) failed: Type: %s, Path: %v, From: %v, To: %v`, changelog.Type, changelog.Path, changelog.From, changelog.To)
			}
		})
	}
}
