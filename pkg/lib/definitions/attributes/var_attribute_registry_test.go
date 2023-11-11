package attributes

import (
	"reflect"
	"testing"
)

func TestLookup(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want *Attribute
	}{
		{"test#1", "name", &name},
		{"test#2", "dn", &distinguishedName},
		{"test#3", "invalid", nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Lookup(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`Lookup(%q) failed: got: %p, want: %p`, tt.args, got, tt.want)
			}
		})
	}
}

func TestLookupMany(t *testing.T) {
	copyRegistry := make(Attributes, len(registry))
	_ = copy(copyRegistry, registry)
	copyRegistry.Sort()

	type args struct {
		strict bool
		attrs  []string
	}

	for _, tt := range []struct {
		name string
		args args
		want Attributes
	}{
		{"test#1", args{true, []string{"name", "name", "accountExpires"}}, Attributes{accountExpires, name}},
		{"test#2", args{true, []string{"invalid"}}, nil},
		{"test#3", args{true, []string{"unicodePwd"}}, Attributes{unicodePassword}},
		{"test#4", args{true, []string{"*"}}, copyRegistry},
		{"test#5", args{false, []string{"*"}}, Attributes{Any()}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := LookupMany(tt.args.strict, tt.args.attrs...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`LookupMany(%v) failed: got: %v, want: %v`, tt.args, got, tt.want)
			}
		})
	}
}

func TestRegistryNotEmpty(t *testing.T) {
	if len(registry) == 0 {
		t.Errorf(`no attributes registered`)
	}
}
