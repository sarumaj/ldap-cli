package auth

import (
	"reflect"
	"testing"

	defaults "github.com/creasty/defaults"
	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

func TestBindParametersToAndFromKeyring(t *testing.T) {
	libutil.SkipOAT(t)

	for _, tt := range []struct {
		name string
		args *BindParameters
		want *BindParameters
	}{
		{"test#1",
			NewBindParameters().SetType(SIMPLE).SetUser("domain\\user").SetPassword("pass"),
			NewBindParameters().SetType(SIMPLE).SetUser("domain\\user").SetPassword("pass")},
		{"test#2",
			NewBindParameters().SetType(NTLM).SetDomain("example.com").SetUser("user").SetPassword("pass"),
			NewBindParameters().SetType(NTLM).SetDomain("example.com").SetUser("user").SetPassword("pass")},
	} {

		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.ToKeyring(); err != nil {
				t.Errorf(`(%T).ToKeyring() failed: %v`, tt.args, err)
			}

			got := NewBindParameters()
			if err := got.FromKeyring(); err != nil {
				t.Errorf(`(%T).FromKeyring() failed: %v`, tt.want, err)
				return
			}

			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf(`(%T).FromKeyring() failed: got: %#v, want: %#v`, tt.want, *got, *tt.want)
			}
		})
	}
}

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
		args *BindParameters
		want *BindParameters
	}{
		{"test#1",
			NewBindParameters().SetType(AuthType(-1)),
			NewBindParameters().SetType(UNAUTHENTICATED)},
		{"test#2",
			NewBindParameters().SetType(NTLM).SetDomain("example.com").SetUser("user").SetPassword("pass"),
			NewBindParameters().SetType(NTLM).SetDomain("example.com").SetUser("user").SetPassword("pass")},
		{"test#3",
			NewBindParameters().SetType(SIMPLE).SetUser("domain\\\\user").SetPassword("pass"),
			NewBindParameters().SetType(SIMPLE).SetUser("domain\\user").SetPassword("pass")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := defaults.Set(tt.args)
			if err != nil {
				t.Errorf(`defaults.Set(%v) failed: %v`, tt.args, err)
			} else if !reflect.DeepEqual(*tt.args, *tt.want) {
				t.Errorf(`defaults.Set(%v) failed: did not expect %v`, *tt.args, *tt.want)
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
