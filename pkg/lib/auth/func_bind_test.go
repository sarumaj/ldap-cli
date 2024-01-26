package auth

import (
	"reflect"
	"testing"

	defaults "github.com/creasty/defaults"
	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

func TestBind(t *testing.T) {
	libutil.SkipOAT(t)

	type args struct {
		*DialOptions
		*BindParameters
	}

	for _, tt := range []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test#1", args{
			NewDialOptions().SetURL("ldap://localhost:389"),
			NewBindParameters().SetType(SIMPLE).SetUser("cn=admin,dc=mock,dc=ad,dc=com").SetPassword("admin"),
		}, false},
		{"test#2", args{
			func() *DialOptions { d := NewDialOptions(); d.SetDefaults(); return d }(),
			NewBindParameters().SetType(UNAUTHENTICATED),
		}, false},
	} {

		t.Run(tt.name, func(t *testing.T) {
			conn, err := Bind(tt.args.BindParameters, tt.args.DialOptions)
			if (err == nil) == tt.wantErr {
				t.Errorf(`Bind(%v, %v) failed: %v`, tt.args.BindParameters, tt.args.DialOptions, err)
				return
			}

			if err == nil && conn != nil {
				if got := conn.RemoteHost(); len(got) == 0 {
					t.Errorf(`conn.RemoteAddr() failed`)
				}

				_ = conn.Close()
			}
		})
	}
}

func TestBindParameterDefaults(t *testing.T) {
	for _, tt := range []struct {
		name string
		args BindParameters
		want BindParameters
	}{
		{"test#1",
			BindParameters{},
			BindParameters{UNAUTHENTICATED, "", "", ""}},
		{"test#2",
			BindParameters{NTLM, "example.com", "user", "pass"},
			BindParameters{NTLM, "example.com", "user", "pass"}},
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
		args    BindParameters
		wantErr bool
	}{
		{"test#1",
			BindParameters{},
			true},
		{"test#2",
			BindParameters{SIMPLE, "", "", ""},
			true},
		{"test#3",
			BindParameters{SIMPLE, "", "user", ""},
			true},
		{"test#4",
			BindParameters{NTLM, "", "user", "pass"},
			true},
		{"test#5",
			BindParameters{NTLM, "example.com", "user", "pass"},
			false},
		{"test#6",
			BindParameters{SIMPLE, "", "user", "pass"},
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
