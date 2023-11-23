package objects

import (
	"fmt"
	"testing"

	"github.com/r3labs/diff/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestGroupRead(t *testing.T) {
	defer util.PatchForTimeNow().Unpatch()
	for _, tt := range []struct {
		name string
		args map[string]any
		want Group
	}{
		{"test#1",
			map[string]any{
				"cn":                "test",
				"description":       "test",
				"displayName":       "test",
				"distinguishedName": "CN=test,...DN:example,DN=com",
				"groupType":         fmt.Sprintf("%d", attributes.GROUP_TYPE_CREATED_BY_SYSTEM|attributes.GROUP_TYPE_GLOBAL|attributes.GROUP_TYPE_SECURITY),
				"memberOf":          []string{"CN=parent,...DN:example,DN=com"},
				"member":            []string{"CN=child,...DN:example,DN=com"},
				"name":              "test",
				"objectCategory":    []string{"group"},
				"objectClass":       []string{"group"},
				"objectGUID":        "696768d7-5ca1-97a6-4835-cbab7c9e5b11",
				"objectSID":         "S-1-5-21-1234567890-123456789-1234567890-1234567",
				"SAMAccountName":    "test",
				"SAMAccountType":    fmt.Sprintf("%d", attributes.SAM_ACCOUNT_TYPE_GROUP_OBJECT),
			},
			Group{
				CommonName:        "test",
				Description:       "test",
				DisplayName:       "test",
				DistinguishedName: "CN=test,...DN:example,DN=com",
				GroupTypeRaw:      2147483651,
				GroupType:         []string{"CREATED_BY_SYSTEM", "GLOBAL", "SECURITY"},
				MemberOf:          []string{"CN=parent,...DN:example,DN=com"},
				Members:           []string{"CN=child,...DN:example,DN=com"},
				Name:              "test",
				ObjectCategory:    "group",
				ObjectClass:       []string{"group"},
				ObjectGUID:        "696768d7-5ca1-97a6-4835-cbab7c9e5b11",
				SamAccountName:    "test",
				SamAccountType:    []string{"GROUP_OBJECT"},
				SamAccountTypeRaw: 268435456,
				SID:               "S-1-5-21-1234567890-123456789-1234567890-1234567"}},
	} {
		var got Group
		if err := got.Read(tt.args); err != nil {
			t.Error(err)
		}

		changelogs, err := diff.Diff(got, tt.want, diff.DisableStructValues())
		if err != nil {
			t.Error(err)
		}

		for _, changelog := range changelogs {
			t.Errorf(`(Group).Read(...) failed: Type: %s, Path: %v, Got: %v, Want: %v`, changelog.Type, changelog.Path, changelog.From, changelog.To)
		}
	}
}
