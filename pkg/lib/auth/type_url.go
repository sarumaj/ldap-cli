package auth

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

var validURLRegex = regexp.MustCompile(`^(?P<Scheme>[^:]+)://(?P<Host>[^:]+):(?P<Port>\d+)`)

// Server's URL
type URL struct {
	// Scheme
	Scheme Scheme `validate:"required,is_valid"`
	// Server's domain name
	Host string `validate:"required"`
	// Server's port
	Port Port `validate:"required,gt=0"`
}

// Get server's hostname and port in form <hostname>:<port>
func (u URL) HostPort() string { return fmt.Sprintf("%s:%d", u.Host, u.Port) }

// Set scheme
func (u *URL) SetScheme(s Scheme) *URL { u.Scheme = s; return u }

// Set hostname
func (u *URL) SetHost(h string) *URL { u.Host = h; return u }

// Set port
func (u *URL) SetPort(p Port) *URL { u.Port = p; return u }

// Render URL as <scheme>://<hostname>:<port>
func (u URL) String() string { return fmt.Sprintf("%s://%s:%d", u.Scheme, u.Host, u.Port) }

// Validate fields
func (u *URL) Validate() error { return util.FormatError(validate.Struct(u)) }

// Make empty URL
func NewURL() *URL { return &URL{} }

// Parse URL from string matching <scheme>://<hostname>:<port>
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

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}
