package auth

import (
	"reflect"
	"testing"

	"github.com/creasty/defaults"
)

func TestBindParameterDefaults(t *testing.T) {
	for _, tt := range []struct {
		name string
		args BindParameter
		want BindParameter
	}{
		{"test#1",
			BindParameter{},
			BindParameter{SIMPLE, "", "", ""}},
		{"test#2",
			BindParameter{NTLM, "example.com", "user", "pass"},
			BindParameter{NTLM, "example.com", "user", "pass"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := &tt.args
			err := defaults.Set(opts)
			if err != nil {
				t.Errorf(`defaults.Set(%v) failed: %v`, tt.args, err)
			} else if !reflect.DeepEqual(*opts, tt.want) {
				t.Errorf(`defaults.Set(%v) failed: did not expect %v`, tt.args, *opts)
			}
		})
	}
}

func TestBindParameterValidate(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    BindParameter
		wantErr bool
	}{
		{"test#1",
			BindParameter{},
			true},
		{"test#2",
			BindParameter{SIMPLE, "", "", ""},
			true},
		{"test#3",
			BindParameter{SIMPLE, "", "user", ""},
			true},
		{"test#4",
			BindParameter{NTLM, "", "user", "pass"},
			true},
		{"test#5",
			BindParameter{NTLM, "example.com", "user", "pass"},
			false},
		{"test#6",
			BindParameter{SIMPLE, "", "user", "pass"},
			false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := &tt.args
			err := opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf(`(%v).Validate() failed: %v`, tt.args, err)
			}
		})
	}
}
