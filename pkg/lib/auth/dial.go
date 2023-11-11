package auth

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/creasty/defaults"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type DialOptions struct {
	Host       string `validate:"required"`
	Port       uint   `validate:"required"`
	SizeLimit  int
	TimeLimit  time.Duration
	MaxRetries uint
	UseTLS     bool
}

func (o *DialOptions) Default() error {
	return defaults.Set(o)
}

func (o *DialOptions) SetDefaults() {
	o.SizeLimit = 10
	o.TimeLimit = 10 * time.Second
	o.MaxRetries = 3
	o.UseTLS = true
}

func (o *DialOptions) Validate() error { return util.FormatError(validate.Struct(o)) }

func Dial(opts *DialOptions) (net.Conn, error) {
	if err := defaults.Set(opts); err != nil {
		opts.MaxRetries = 0
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		opts.MaxRetries = 0
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	if opts.UseTLS {
		return tls.DialWithDialer(
			&net.Dialer{
				Timeout: time.Duration(opts.TimeLimit) * time.Second,
			},
			"tcp",
			addr,
			&tls.Config{},
		)
	}

	return net.DialTimeout(
		"tcp",
		addr,
		time.Duration(opts.TimeLimit)*time.Second,
	)
}
