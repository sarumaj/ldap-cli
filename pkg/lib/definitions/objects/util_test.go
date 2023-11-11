package objects

import (
	"testing"

	diff "github.com/r3labs/diff/v3"
)

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
	// TODO
	type object struct {
	}
	type args struct {
		o object
		m map[string]any
	}
	for _, tt := range []struct {
		name string
		args args
		want object
	}{} {
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
