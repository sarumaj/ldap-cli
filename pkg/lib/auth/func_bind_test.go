package auth

import (
	"crypto/tls"
	"os"
	"reflect"
	"testing"

	"github.com/creasty/defaults"
)

func TestBind(t *testing.T) {
	type args struct {
		*DialOptions
		*BindParameters
	}

	for _, tt := range []struct {
		name string
		args args
	}{
		{"test#1", args{
			NewDialOptions().SetURL(os.Getenv("AD_AUTO_URL")).SetTLSConfig(&tls.Config{InsecureSkipVerify: true}),
			NewBindParameters().SetType(SIMPLE).SetDomain("").SetUser(os.Getenv("AD_DEFAULT_USER")).SetPassword(os.Getenv("AD_DEFAULT_PASS")),
		}},
		{"test#2", args{
			NewDialOptions().SetURL(os.Getenv("AD_DMZ01_URL")).SetTLSConfig(&tls.Config{InsecureSkipVerify: true}),
			nil,
		}},
	} {

		t.Run(tt.name, func(t *testing.T) {
			if tt.args.URL == nil {
				t.Skipf("Skipping %s", tt.name)
			}

			conn, err := Bind(tt.args.BindParameters, tt.args.DialOptions)
			if err != nil {
				t.Errorf(`Bind(%v, %v) failed: %v`, tt.args.BindParameters, tt.args.DialOptions, err)
				return
			}
			_ = conn.Close()
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
