package auth

import (
	"crypto/tls"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/creasty/defaults"
)

func TestDial(t *testing.T) {
	for _, tt := range []struct {
		name string
		args *DialOptions
		skip bool
	}{
		{"test#1", (&DialOptions{}).SetURL("ldaps://auto.contiwan.com:636").SetTLSConfig(&tls.Config{InsecureSkipVerify: true}), false},
		{"test#2", (&DialOptions{}).SetURL("ldaps://dmz01.net:636").SetTLSConfig(&tls.Config{InsecureSkipVerify: true}), false},
	} {

		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf("Skipping %s", tt.name)
			}

			_, err := Dial(tt.args)
			if err != nil {
				t.Errorf(`Dial(%v) failed: %v`, tt.args, err)
			}
		})
	}
}

func TestDialOptionsDefaults(t *testing.T) {
	for _, tt := range []struct {
		name string
		args DialOptions
		want DialOptions
	}{
		{"test#1",
			DialOptions{0, 0, 0, nil, nil},
			DialOptions{3, 10, 10 * time.Second, nil, &url.URL{Scheme: "ldap", Host: "localhost:389"}}},
		{"test#2",
			DialOptions{5, 20, time.Second, &tls.Config{}, &url.URL{Scheme: "ldaps", Host: "example.com:389"}},
			DialOptions{5, 20, time.Second, &tls.Config{}, &url.URL{Scheme: "ldaps", Host: "example.com:389"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := &tt.args
			err := defaults.Set(opts)
			if err != nil {
				t.Errorf(`defaults.Set(DialOptions) failed: %v`, err)
			} else if !reflect.DeepEqual(*opts, tt.want) {
				t.Errorf(`defaults.Set(DialOptions) failed: did not expect %v`, *opts)
			}

		})
	}
}

func TestDialOptionsVerify(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    DialOptions
		wantErr bool
	}{
		{"test#1",
			DialOptions{0, 0, 0, nil, nil},
			true},
		{"test#2",
			DialOptions{5, 20, time.Second, &tls.Config{}, &url.URL{Scheme: "ldaps", Host: "example.com:389"}},
			false},
		{"test#3",
			DialOptions{5, 20, time.Second, nil, &url.URL{Scheme: "ldaps", Host: "example.com:389"}},
			false},
		{"test#4",
			DialOptions{5, 20, time.Second, nil, nil},
			true},
		{"test#5",
			DialOptions{5, 20, 0, nil, &url.URL{Scheme: "ldaps", Host: "example.com:389"}},
			true},
		{"test#6",
			DialOptions{5, 0, time.Second, nil, &url.URL{Scheme: "ldaps", Host: "example.com:389"}},
			true},
		{"test#7",
			DialOptions{0, 20, time.Second, nil, &url.URL{Scheme: "ldaps", Host: "example.com:389"}},
			true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := &tt.args
			err := opts.Validate()
			if err != nil && !tt.wantErr {
				t.Errorf(`(DialOptions).Validate() failed: %v`, err)
			}
		})
	}
}
