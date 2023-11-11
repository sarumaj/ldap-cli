package auth

import (
	"fmt"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type URL struct {
	Scheme Scheme `validate:"required,is_valid"`
	Host   string `validate:"required"`
	Port   Port   `validate:"required"`
}

func (u URL) HostPort() string         { return fmt.Sprintf("%s:%d", u.Host, u.Port) }
func (u *URL) SetScheme(s Scheme) *URL { u.Scheme = s; return u }
func (u *URL) SetHost(h string) *URL   { u.Host = h; return u }
func (u *URL) SetPort(p Port) *URL     { u.Port = p; return u }
func (u URL) String() string           { return fmt.Sprintf("%s://%s:%d", u.Scheme, u.Host, u.Port) }
func (u *URL) Validate() error         { return util.FormatError(validate.Struct(u)) }

func URLFromString(in string) (*URL, error) {
	var u URL
	if _, err := fmt.Sscanf(in, "%s://%s:%d", &u.Scheme, &u.Host, &u.Port); err != nil {
		return nil, err
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return &u, nil
}
