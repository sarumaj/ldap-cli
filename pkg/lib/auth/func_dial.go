package auth

import (
	"crypto/tls"
	"math/rand"
	"net"
	"time"

	defaults "github.com/creasty/defaults"
	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

// Options for dialer
type DialOptions struct {
	// Number of max retries if failing
	MaxRetries uint `validate:"required" default:"3"`
	// Limits number of objects returned by an LDAP query
	SizeLimit int64
	// Timeout for connection handshake and LDAP queries
	TimeLimit time.Duration `validate:"required" default:"10s"`
	// Custom TLS config
	TLSConfig *tls.Config
	// Server URL
	URL *URL `validate:"required,is_valid"` // default: ldap://localhost:389
}

// Sets default URL
func (o *DialOptions) SetDefaults() {
	if o.URL == nil {
		o.URL = &URL{
			Scheme: LDAP,
			Host:   "localhost",
			Port:   LDAP_RW,
		}
	}
}

// SetMaxRetries sets max retries
func (o *DialOptions) SetMaxRetries(retries uint) *DialOptions { o.MaxRetries = retries; return o }

// SetSizeLimit sets size limit
func (o *DialOptions) SetSizeLimit(limit int64) *DialOptions { o.SizeLimit = limit; return o }

// SetTimeLimit sets time limit
func (o *DialOptions) SetTimeLimit(limit time.Duration) *DialOptions { o.TimeLimit = limit; return o }

// SetURL sets URL
func (o *DialOptions) SetURL(addr string) *DialOptions {
	o.URL, _ = URLFromString(addr)
	return o
}

// SetTLSConfig sets TLS config
func (o *DialOptions) SetTLSConfig(conf *tls.Config) *DialOptions { o.TLSConfig = conf; return o }

// Validate validates options
func (o *DialOptions) Validate() error { return libutil.FormatError(validate.Struct(o)) }

// Dial connects to an LDAP server
func Dial(opts *DialOptions) (conn net.Conn, err error) {
	if opts == nil {
		opts = NewDialOptions()
	}

	if err := defaults.Set(opts); err != nil {
		opts.MaxRetries = 0 // abort immediately
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		opts.MaxRetries = 0 // abort immediately
		return nil, err
	}

	if opts.URL.Scheme == LDAPI {
		conn, err = net.DialTimeout("unix", "/var/run/slapd/ldapi", opts.TimeLimit)
		return conn, libutil.Handle(err)
	}

	switch opts.URL.Host {
	case "localhost": // ignore

	default:
		// resolve DNS
		addresses, err := net.LookupHost(opts.URL.Host)
		if err != nil {
			return nil, err
		}

		// select random address and resolve back
		opts.URL.Host = libutil.LookupAddress(addresses[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(addresses))])

	}

	// switch over to TLS if necessary
	switch opts.URL.Scheme {
	case LDAPS:
		if opts.TLSConfig == nil {
			opts.TLSConfig = &tls.Config{}
		}

		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: opts.TimeLimit}, "tcp", opts.URL.HostPort(), opts.TLSConfig)

	case CLDAP:
		conn, err = net.DialTimeout("udp", opts.URL.HostPort(), opts.TimeLimit)

	default:
		conn, err = net.DialTimeout("tcp", opts.URL.HostPort(), opts.TimeLimit)

	}

	return conn, libutil.Handle(err)
}

// NewDialOptions creates new options
func NewDialOptions() *DialOptions { return &DialOptions{} }
