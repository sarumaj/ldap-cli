package auth

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type DialOptions struct {
	MaxRetries uint          `validate:"required" default:"3"`
	SizeLimit  int64         `validate:"required" default:"10"`
	TimeLimit  time.Duration `validate:"required" default:"10s"`
	TLSConfig  *tls.Config
	URL        *url.URL `validate:"required"`
}

func (o *DialOptions) SetDefaults() {
	if o.URL == nil {
		o.URL, _ = url.Parse(fmt.Sprintf("ldap://localhost:%d", LDAP))
	}
}

func (o *DialOptions) SetMaxRetries(retries uint) *DialOptions       { o.MaxRetries = retries; return o }
func (o *DialOptions) SetSizeLimit(limit int64) *DialOptions         { o.SizeLimit = limit; return o }
func (o *DialOptions) SetTimeLimit(limit time.Duration) *DialOptions { o.TimeLimit = limit; return o }

func (o *DialOptions) SetURL(addr string) *DialOptions {
	parsedURL, err := url.Parse(addr)
	if err != nil {
		return o
	}

	o.URL = parsedURL
	return o
}

func (o *DialOptions) SetTLSConfig(conf *tls.Config) *DialOptions { o.TLSConfig = conf; return o }
func (o *DialOptions) Validate() error                            { return util.FormatError(validate.Struct(o)) }

func Dial(opts *DialOptions) (net.Conn, error) {
	if err := defaults.Set(opts); err != nil {
		opts.MaxRetries = 0 // abort immediately
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		opts.MaxRetries = 0 // abort immediately
		return nil, err
	}

	if strings.HasSuffix(opts.URL.Scheme, "s") {
		if opts.TLSConfig == nil {
			opts.TLSConfig = &tls.Config{}
		}

		return tls.DialWithDialer(&net.Dialer{Timeout: opts.TimeLimit}, "tcp", opts.URL.Host, opts.TLSConfig)
	}

	return net.DialTimeout("tcp", opts.URL.Host, opts.TimeLimit)
}
