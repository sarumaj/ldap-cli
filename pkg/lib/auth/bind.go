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

type BindParameter struct {
	Type     Type   `validate:"required,is_valid"`
	Domain   string `validate:"required_if=Type NTLM"`
	Password string `validate:"required"`
	User     string `validate:"required"`
}

func (p *BindParameter) SetDefaults()    { p.Type = SIMPLE }
func (p *BindParameter) Validate() error { return util.FormatError(validate.Struct(p)) }

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

	conn.remoteHost = c.RemoteAddr().String()
	fmt.Sprintln(conn.remoteHost)
	raw := strings.Split(conn.remoteHost, ":")
	if addr, err := net.LookupAddr(raw[0]); err == nil && len(addr) > 0 {
		fmt.Sprintln(addr)
		conn.remoteHost = fmt.Sprintf("%s:%d", strings.Trim(addr[0], "."), conn.options.Port)
	}

	ldapConn := ldap.NewConn(c, true)
	ldapConn.SetTimeout(conn.options.TimeLimit)
	ldapConn.Start()

	switch parameters.Type {

	case SIMPLE:
		err = ldapConn.Bind(parameters.User, parameters.Password)

	case MD5:
		err = ldapConn.MD5Bind(conn.options.Host, parameters.User, parameters.Password)

	case NTLM:
		err = ldapConn.NTLMBind(parameters.Domain, parameters.User, parameters.Password)

	}
	if err != nil {
		return nil, err
	}

	conn.connection = ldapConn
	return conn, nil
}
