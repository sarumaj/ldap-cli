package auth

import (
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-ldap/ldap/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

// Bind parameters
type BindParameters struct {
	// Type of authentication
	AuthType AuthType `validate:"required,is_valid"` // default: SIMPLE
	// User's domain (required for NTLM authentication)
	Domain string `validate:"required_if=AuthType NTLM"`
	// User's password
	Password string `validate:"required_unless=AuthType UNAUTHENTICATED"`
	// Username
	User string `validate:"required_unless=AuthType UNAUTHENTICATED"`
}

// Set default Type
func (p *BindParameters) SetDefaults() {
	if p.AuthType == 0 || !p.AuthType.IsValid() {
		p.AuthType = UNAUTHENTICATED
	}

	if i := strings.Index(p.User, `\\`); i > 0 {
		p.User = strings.Replace(p.User, `\\`, `\`, 1)
	}
}

// Set domain (required for NTLM-based authentication)
func (p *BindParameters) SetDomain(domain string) *BindParameters {
	p.Domain = domain
	return p
}

// Set password
func (p *BindParameters) SetPassword(password string) *BindParameters {
	p.Password = password
	return p
}

// Set username
func (p *BindParameters) SetUser(user string) *BindParameters {
	p.User = user
	return p
}

// Set authentication type
func (p *BindParameters) SetType(authType AuthType) *BindParameters {
	p.AuthType = authType
	return p
}

// Validate fields
func (p *BindParameters) Validate() error { return util.FormatError(validate.Struct(p)) }

// Establish connection with the server
func Bind(parameters *BindParameters, options *DialOptions) (*Connection, error) {
	if parameters == nil {
		parameters = NewBindParameters()
	}

	if err := defaults.Set(parameters); err != nil {
		return nil, err
	}

	if err := parameters.Validate(); err != nil {
		return nil, err
	}

	c, err := Dial(options)
	for i, d := uint(0), time.Second; i < options.MaxRetries && err != nil; i, d = i+1, d*2 {
		<-time.After(d)
		c, err = Dial(options)
	}
	if err != nil {
		return nil, err
	}

	ldapConn := ldap.NewConn(c, true)
	ldapConn.SetTimeout(options.TimeLimit)
	ldapConn.Start()

	switch parameters.AuthType {

	case UNAUTHENTICATED:
		err = ldapConn.UnauthenticatedBind(parameters.User)

	case SIMPLE:
		err = ldapConn.Bind(parameters.User, parameters.Password)

	case MD5:
		err = ldapConn.MD5Bind(options.URL.Host, parameters.User, parameters.Password)

	case NTLM:
		err = ldapConn.NTLMBind(parameters.Domain, parameters.User, parameters.Password)

	}
	if err != nil {
		return nil, err
	}

	return &Connection{
		Conn:        ldapConn,
		DialOptions: options,
		remoteHost:  util.LookupAddress(c.RemoteAddr().String()),
	}, nil
}

func NewBindParameters() *BindParameters { return &BindParameters{AuthType: UNAUTHENTICATED} }
