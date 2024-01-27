package client

import (
	"testing"
	"time"

	auth "github.com/sarumaj/ldap-cli/v2/pkg/lib/auth"
	attributes "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
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

	t.Cleanup(func() {
		_ = conn.Modify(ModifyGroupMembersRequest(
			"cn=group02,dc=mock,dc=ad,dc=com",
			nil, []string{"uix00002"}, nil,
			attributes.Attribute{LDAPDisplayName: "memberUid"},
		))
		_ = conn.Close()
	})

	type args struct {
		add, delete, replace []string
	}

	for _, tt := range []struct {
		name string
		args args
	}{
		{"test#1", args{[]string{"uix00002"}, nil, nil}},
		{"test#2", args{nil, []string{"uix00002"}, nil}},
		{"test#3", args{nil, nil, []string{"uix00002"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			request := ModifyGroupMembersRequest(
				"cn=group02,dc=mock,dc=ad,dc=com",
				tt.args.add, tt.args.delete, tt.args.replace,
				attributes.Attribute{LDAPDisplayName: "memberUid"},
			)
			if err := libutil.Handle(conn.Modify(request)); err != nil {
				t.Errorf(
					`ModifyGroupMembersRequest("cn=group02,dc=mock,dc=ad,dc=com", %v, %v, %v, attributes.Attribute{LDAPDisplayName: "memberUid"}) failed: %v`,
					tt.args.add, tt.args.delete, tt.args.replace, err,
				)
			}
		})
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
