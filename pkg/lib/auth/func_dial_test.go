package auth

import (
	"crypto/tls"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/creasty/defaults"
)

func TestDial(t *testing.T) {
	for _, tt := range []struct {
		name string
		args *DialOptions
	}{
		{"test#1", (&DialOptions{}).SetURL(os.Getenv("AD_AUTO_URL")).SetTLSConfig(&tls.Config{InsecureSkipVerify: true})},
		{"test#2", (&DialOptions{}).SetURL(os.Getenv("AD_DMZ01_URL")).SetTLSConfig(&tls.Config{InsecureSkipVerify: true})},
	} {

		t.Run(tt.name, func(t *testing.T) {
			if tt.args.URL == nil {
				t.Skipf("Skipping %s", tt.name)
			}

			conn, err := Dial(tt.args)
			if err != nil {
				t.Errorf(`Dial(%v) failed: %v`, tt.args, err)
				return
			}
			_ = conn.Close()
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
			DialOptions{3, 10, 10 * time.Second, nil, &URL{"ldap", "localhost", 389}}},
		{"test#2",
			DialOptions{5, 20, time.Second, &tls.Config{}, &URL{"ldaps", "example.com", 389}},
			DialOptions{5, 20, time.Second, &tls.Config{}, &URL{"ldaps", "example.com", 389}}},
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

func TestDialOptionsValidate(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    DialOptions
		wantErr bool
	}{
		{"test#1",
			DialOptions{0, 0, 0, nil, nil},
			true},
		{"test#2",
			DialOptions{5, 20, time.Second, &tls.Config{}, &URL{Scheme: "ldaps", Host: "example.com", Port: 389}},
			false},
		{"test#3",
			DialOptions{5, 20, time.Second, nil, &URL{Scheme: "ldaps", Host: "example.com", Port: 389}},
			false},
		{"test#4",
			DialOptions{5, 20, time.Second, nil, nil},
			true},
		{"test#5",
			DialOptions{5, 20, 0, nil, &URL{Scheme: "ldaps", Host: "example.com", Port: 389}},
			true},
		{"test#6",
			DialOptions{5, 0, time.Second, nil, &URL{Scheme: "ldaps", Host: "example.com", Port: 389}},
			true},
		{"test#7",
			DialOptions{0, 20, time.Second, nil, &URL{Scheme: "ldaps", Host: "example.com", Port: 389}},
			true},
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
