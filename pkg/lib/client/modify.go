package client

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/auth"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func modify(conn *auth.Connection, request *ldap.ModifyRequest) error {
	if conn == nil {
		return libutil.ErrBindFirst
	}

	if err := conn.Modify(request); err != nil {
		return err
	}

	return nil
}

func modifyAddMembers(conn *auth.Connection, dn string, members []string) error {
	request := ldap.NewModifyRequest(dn, nil)
	request.Add(attributes.Members().String(), members)

	return modify(conn, request)
}

func modifyRemoveMembers(conn *auth.Connection, dn string, members []string) error {
	request := ldap.NewModifyRequest(dn, nil)
	request.Delete(attributes.Members().String(), members)

	return modify(conn, request)
}

func modifyReplaceMembers(conn *auth.Connection, dn string, members []string) error {
	request := ldap.NewModifyRequest(dn, nil)
	request.Replace(attributes.Members().String(), members)

	return modify(conn, request)
}

func modifyReplacePassword(conn *auth.Connection, dn, oldPassword, newPassword string) error {
	_, err := conn.PasswordModify(ldap.NewPasswordModifyRequest(dn, oldPassword, newPassword))
	return err
}
