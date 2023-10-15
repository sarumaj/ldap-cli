package auth

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-ldap/ldap/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

// Bind parameters
type BindParameter struct {
	// Type of authentication
	Type Type `validate:"required,is_valid"` // default: SIMPLE
	// User's domain (required for NTLM authentication)
	Domain string `validate:"required_if=Type NTLM"`
	// User's password
	Password string `validate:"required"`
	// Username
	User string `validate:"required"`
}

// Set default Type
func (p *BindParameter) SetDefaults() {
	if p.Type == 0 || !p.Type.IsValid() {
		p.Type = SIMPLE
	}
}

// Validate fields
func (p *BindParameter) Validate() error { return util.FormatError(validate.Struct(p)) }

// Establish connection with the server
func Bind(parameters *BindParameter, options *DialOptions) (*Connection, error) {
	if err := defaults.Set(parameters); err != nil {
		return nil, err
	}

	if err := parameters.Validate(); err != nil {
		return nil, err
	}

	conn := &Connection{
		options: options,
	}

	var c net.Conn
	var err error
	for i, d := uint(0), time.Second; i < conn.options.MaxRetries; i, d = i+1, d*2 {
		c, err = Dial(conn.options)
		if err == nil {
			break
		}

		time.Sleep(d)
	}
	if err != nil {
		return nil, err
	}

	// TODO: examine necessity
	conn.remoteHost = c.RemoteAddr().String()
	fmt.Sprintln(conn.remoteHost)
	raw := strings.Split(conn.remoteHost, ":")
	if addr, err := net.LookupAddr(raw[0]); err == nil && len(addr) > 0 {
		fmt.Sprintln(addr)
		conn.remoteHost = fmt.Sprintf("%s:%d", strings.Trim(addr[0], "."), conn.options.URL.Port)
	}

	ldapConn := ldap.NewConn(c, true)
	ldapConn.SetTimeout(conn.options.TimeLimit)
	ldapConn.Start()

	switch parameters.Type {

	case SIMPLE:
		err = ldapConn.Bind(parameters.User, parameters.Password)

	case MD5:
		err = ldapConn.MD5Bind(conn.options.URL.Host, parameters.User, parameters.Password)

	case NTLM:
		err = ldapConn.NTLMBind(parameters.Domain, parameters.User, parameters.Password)

	}
	if err != nil {
		return nil, err
	}

	conn.connection = ldapConn
	return conn, nil
}
