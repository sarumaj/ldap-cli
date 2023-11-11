package client

import (
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

func ModifyGroupMembersRequest(dn string, add, delete, replace []string) *ldap.ModifyRequest {
	request := ldap.NewModifyRequest(dn, nil)
	if len(add) > 0 {
		request.Add(attributes.Members().String(), add)
	}

	if len(delete) > 0 {
		request.Delete(attributes.Members().String(), delete)
	}

	if len(replace) > 0 {
		request.Replace(attributes.Members().String(), replace)
	}

	return request

}

func ModifyPasswordRequest(dn, newPassword string) *ldap.ModifyRequest {
	request := ldap.NewModifyRequest(dn, nil)
	request.Replace(attributes.UserPassword().String(), []string{newPassword})
	return request
}
