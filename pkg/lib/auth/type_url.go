package auth

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

// validURLRegex is a regex to validate URL
var validURLRegex = regexp.MustCompile(`^(?P<Scheme>[^:]+)://(?P<Host>[^:]+)(?::(?P<Port>\d+))?`)

var _ libutil.ValidatorInterface = URL{}

// URL is a URL (RFC 4516)
type URL struct {
	// Scheme
	Scheme Scheme `validate:"required,is_valid"`
	// Server's domain name
	Host string `validate:"required"`
	// Server's port
	Port Port `validate:"excluded_if=Scheme ldapi,required_unless=Scheme ldapi"`
}

// HostPort returns <hostname>:<port>
func (u URL) HostPort() string {
	if u.Port > 0 {
		return fmt.Sprintf("%s:%d", u.Host, u.Port)
	}
	return u.Host
}

// IsValid returns true if the URL is valid
func (u URL) IsValid() bool { return validate.Struct(&u) == nil }

// SetScheme sets scheme
func (u *URL) SetScheme(s Scheme) *URL { u.Scheme = s; return u }

// SetHost sets host
func (u *URL) SetHost(h string) *URL { u.Host = h; return u }

// SetPort sets port
func (u *URL) SetPort(p Port) *URL { u.Port = p; return u }

// String returns the string representation of an URL (<scheme>://<hostname>:<?port>)
func (u URL) String() string {
	if u.Port > 0 {
		return fmt.Sprintf("%s://%s:%d", u.Scheme, u.Host, u.Port)
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

// ToBaseDirectoryPath builds base DN from host
func (u URL) ToBaseDirectoryPath() string {
	var components []string
	for _, dc := range strings.Split(u.Host, ".") {
		if dc == "" {
			continue
		}

		components = append(components, "DC="+dc)
	}

	return strings.Join(components, ",")
}

// Validate validates the URL
func (u *URL) Validate() error { return libutil.FormatError(validate.Struct(u)) }

// NewURL returns a new URL
func NewURL() *URL { return &URL{} }

// URLFromString returns an URL from a string (<scheme>://<hostname>:<port>)
func URLFromString(in string) (*URL, error) {
	if !validURLRegex.MatchString(in) {
		return nil, fmt.Errorf("provided address %q does not match the validation pattern: %q", in, validURLRegex)
	}

	matches := validURLRegex.FindStringSubmatch(in)
	u := &URL{
		Scheme: Scheme(matches[1]),
		Host:   matches[2],
	}

	port, _ := strconv.Atoi(matches[3])
	u.Port = Port(port)

	if u.Port == 0 {
		switch u.Scheme {
		case LDAPS:
			u.Port = LDAPS_RW

		case LDAP, CLDAP:
			u.Port = LDAP_RW

		}
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}
