package client

import (
	"testing"
	"time"

	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestModifyGroupMembersRequest(t *testing.T) {
	libutil.SkipOAT(t)

	conn, err := auth.Bind(
		auth.NewBindParameters().SetType(auth.SIMPLE).SetUser("cn=admin,dc=mock,dc=ad,dc=com").SetPassword("admin"),
		auth.NewDialOptions().SetSizeLimit(10).SetTimeLimit(time.Minute*5),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	memberUIDs := []string{"uix00002"}
	request := ModifyGroupMembersRequest("cn=group02,dc=mock,dc=ad,dc=com", memberUIDs, nil, nil, attributes.Attribute{LDAPDisplayName: "memberUid"})
	if err := conn.Modify(request); err != nil {
		t.Error(err)
	}

	request = ModifyGroupMembersRequest("cn=group02,dc=mock,dc=ad,dc=com", nil, memberUIDs, nil, attributes.Attribute{LDAPDisplayName: "memberUid"})
	if err := conn.Modify(request); err != nil {
		t.Error(err)
	}
}

func TestModifyPasswordRequest(t *testing.T) {
	libutil.SkipOAT(t)

	conn, err := auth.Bind(
		auth.NewBindParameters().SetType(auth.SIMPLE).SetUser("cn=admin,dc=mock,dc=ad,dc=com").SetPassword("admin"),
		auth.NewDialOptions().SetSizeLimit(10).SetTimeLimit(time.Minute*5),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	request := ModifyPasswordRequest("uid=uix00002,dc=mock,dc=ad,dc=com", "987654321", attributes.Attribute{LDAPDisplayName: "userPassword"})
	if err := conn.Modify(request); err != nil {
		t.Error(err)
	}
}
