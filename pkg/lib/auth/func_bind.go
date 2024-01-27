package auth

import (
	"strings"
	"time"

	defaults "github.com/creasty/defaults"
	ldap "github.com/go-ldap/ldap/v3"
	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

// BindParameters are parameters for binding to the server
type BindParameters struct {
	// AuthType is the authentication type
	AuthType AuthType `validate:"required,is_valid"` // default: SIMPLE
	// Domain is user's domain (required for NTLM authentication)
	Domain string `validate:"required_if=AuthType NTLM"`
	// User's password
	Password string `validate:"required_unless=AuthType UNAUTHENTICATED"`
	// Username
	User string `validate:"required_unless=AuthType UNAUTHENTICATED"`
}

// FromKeyring loads credentials from keyring
func (p *BindParameters) FromKeyring() error {
	var err error
	if p.User == "" {
		p.User, err = libutil.GetFromKeyring("user")
		if err != nil {
			return err
		}
	}

	if p.Password == "" {
		p.Password, err = libutil.GetFromKeyring("password")
		if err != nil {
			return err
		}
	}

	if p.Domain == "" {
		p.Domain, err = libutil.GetFromKeyring("domain")
		if err != nil {
			return err
		}
	}

	if p.AuthType <= UNAUTHENTICATED || !p.AuthType.IsValid() {
		authType, err := libutil.GetFromKeyring("type")
		if err != nil {
			return err
		}

		p.AuthType = TypeFromString(authType)
	}

	return nil
}

// SetDefaults sets default values
func (p *BindParameters) SetDefaults() {
	if p.AuthType == 0 || !p.AuthType.IsValid() {
		p.AuthType = UNAUTHENTICATED
	}

	if i := strings.Index(p.User, `\\`); i > 0 {
		p.User = strings.Replace(p.User, `\\`, `\`, 1)
	}
}

// SetDomain sets domain (required for NTLM-based authentication)
func (p *BindParameters) SetDomain(domain string) *BindParameters {
	p.Domain = domain
	return p
}

// SetPassword sets password
func (p *BindParameters) SetPassword(password string) *BindParameters {
	p.Password = password
	return p
}

// ToKeyring saves credentials to keyring
func (p BindParameters) ToKeyring() error {
	if err := libutil.SetToKeyring("user", p.User); err != nil {
		return err
	}

	if err := libutil.SetToKeyring("password", p.Password); err != nil {
		return err
	}

	if err := libutil.SetToKeyring("domain", p.Domain); err != nil {
		return err
	}

	if err := libutil.SetToKeyring("type", p.AuthType.String()); err != nil {
		return err
	}

	return nil
}

// SetUser sets username
func (p *BindParameters) SetUser(user string) *BindParameters {
	p.User = user
	return p
}

// SetType sets authentication type
func (p *BindParameters) SetType(authType AuthType) *BindParameters {
	p.AuthType = authType
	return p
}

// Validate validates bind parameters
func (p *BindParameters) Validate() error { return libutil.FormatError(validate.Struct(p)) }

// Bind establishes a connection to the server and binds to it
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
		return nil, libutil.Handle(err)
	}

	ldapConn := ldap.NewConn(c, true)
	ldapConn.SetTimeout(options.TimeLimit)
	ldapConn.Start()

	switch parameters.AuthType {

	case UNAUTHENTICATED:
		err = libutil.Handle(ldapConn.UnauthenticatedBind(parameters.User))

	case SIMPLE:
		err = libutil.Handle(ldapConn.Bind(parameters.User, parameters.Password))

	case MD5:
		err = libutil.Handle(ldapConn.MD5Bind(options.URL.Host, parameters.User, parameters.Password))

	case NTLM:
		err = libutil.Handle(ldapConn.NTLMBind(parameters.Domain, parameters.User, parameters.Password))

	}
	if err != nil {
		return nil, err
	}

	return &Connection{
		Conn:        ldapConn,
		DialOptions: options,
		remoteHost:  libutil.LookupAddress(c.RemoteAddr().String()),
	}, nil
}

// NewBindParameters creates a new BindParameters instance
func NewBindParameters() *BindParameters { return &BindParameters{AuthType: UNAUTHENTICATED} }
