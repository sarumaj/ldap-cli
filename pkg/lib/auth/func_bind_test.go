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
		*BindParameter
	}

	for _, tt := range []struct {
		name string
		args args
	}{
		{"test#1", args{
			(&DialOptions{}).SetURL(os.Getenv("AD_AUTO_URL")).SetTLSConfig(&tls.Config{InsecureSkipVerify: true}),
			(&BindParameter{}).SetType(SIMPLE).SetDomain("").SetUser(os.Getenv("AD_DEFAULT_USER")).SetPassword(os.Getenv("AD_DEFAULT_PASS")),
		}},
		{"test#2", args{
			(&DialOptions{}).SetURL(os.Getenv("AD_DMZ01_URL")).SetTLSConfig(&tls.Config{InsecureSkipVerify: true}),
			nil,
		}},
	} {

		t.Run(tt.name, func(t *testing.T) {
			if tt.args.URL == nil {
				t.Skipf("Skipping %s", tt.name)
			}

			conn, err := Bind(tt.args.BindParameter, tt.args.DialOptions)
			if err != nil {
				t.Errorf(`Bind(%v, %v) failed: %v`, tt.args.BindParameter, tt.args.DialOptions, err)
				return
			}
			_ = conn.Close()
		})
	}
}

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
