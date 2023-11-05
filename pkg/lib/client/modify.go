package client

import (
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

func ModifyGroupMembersRequest(dn string, add, delete, replace []string, memberAttribute attributes.Attribute) *ldap.ModifyRequest {
	request := ldap.NewModifyRequest(dn, nil)
	if len(add) > 0 {
		request.Add(memberAttribute.String(), add)
	}

	if len(delete) > 0 {
		request.Delete(memberAttribute.String(), delete)
	}

	if len(replace) > 0 {
		request.Replace(memberAttribute.String(), replace)
	}

	return request

}

func ModifyPasswordRequest(dn, newPassword string, passwordAttribute attributes.Attribute) *ldap.ModifyRequest {
	request := ldap.NewModifyRequest(dn, nil)
	request.Replace(passwordAttribute.String(), []string{newPassword})
	return request
}
