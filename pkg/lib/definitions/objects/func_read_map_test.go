package objects

import (
	"testing"

	"github.com/r3labs/diff/v3"
)

func TestReadMap(t *testing.T) {
	type object struct {
		Strings []string `ldap_attr:"stringList"`
		Bool    bool
		Int     int
		Int64   int64
		String  string
		Ignore  any `ldap_attr:"-"`
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
				"stringList": "test",
				"bool":       "true",
				"int":        "12",
				"int64":      "13",
				"string":     []string{"test#1.1", "test#1.2"},
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
				"stringList": []string{"test#1.1", "test#1.2"},
				"bool":       "true",
				"int":        "12",
				"int64":      "13",
				"string":     "test",
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
			if err := readMap(&got, tt.args.m); err != nil {
				t.Error(err)
			}

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
