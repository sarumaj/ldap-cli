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
	Host       string        `validate:"required"`
	Port       uint          `validate:"required"`
	SizeLimit  int           `validate:"required" default:"10"`
	TimeLimit  time.Duration `validate:"required" default:"10s"`
	MaxRetries uint          `validate:"required" default:"3"`
	TLSConfig  *tls.Config   `validate:"required" default:"{}"`
}

func (o *DialOptions) Default() error {
	return defaults.Set(o)
}

func (o *DialOptions) Validate() error { return util.FormatError(validate.Struct(o)) }

func Dial(opts *DialOptions) (net.Conn, error) {
	if err := defaults.Set(opts); err != nil {
		opts.MaxRetries = 0 // abort immediately
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		opts.MaxRetries = 0 // abort immediately
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	if opts.TLSConfig != nil {
		return tls.DialWithDialer(&net.Dialer{Timeout: opts.TimeLimit}, "tcp", addr, opts.TLSConfig)
	}

	return net.DialTimeout("tcp", addr, opts.TimeLimit)
}
