package objects

import (
	"fmt"
	"testing"
	"time"

	"github.com/r3labs/diff/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestDomainControllerRead(t *testing.T) {
	defer util.PatchForTimeNow().Unpatch()
	for _, tt := range []struct {
		name string
		args map[string]any
		want DomainController
	}{
		{"test#1",
			map[string]any{
				"accountExpires":     "1234567891234567811",
				"description":        "test",
				"distinguishedName":  "CN=test,...DN:example,DN=com",
				"dnsHostname":        "test",
				"objectCategory":     []string{"computer"},
				"objectClass":        []string{"computer"},
				"objectGUID":         "696768d7-5ca1-97a6-4835-cbab7c9e5b11",
				"objectSID":          "S-1-5-21-1234567890-123456789-1234567890-1234567",
				"SAMAccountName":     "test",
				"SAMAccountType":     fmt.Sprintf("%d", attributes.SAM_ACCOUNT_TYPE_DOMAIN_OBJECT),
				"userAccountControl": fmt.Sprintf("%d", attributes.USER_ACCOUNT_CONTROL_SERVER_TRUST_ACCOUNT),
			},
			DomainController{
				AccountExpires:        1234567891234567811,
				AccountExpiryDate:     time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC),
				Description:           "test",
				DistinguishedName:     "CN=test,...DN:example,DN=com",
				Enabled:               true,
				Hostname:              "test",
				ObjectCategory:        "computer",
				ObjectClass:           []string{"computer"},
				ObjectGUID:            "696768d7-5ca1-97a6-4835-cbab7c9e5b11",
				SamAccountName:        "test",
				SamAccountType:        []string{"DOMAIN_OBJECT"},
				SamAccountTypeRaw:     0,
				SID:                   "S-1-5-21-1234567890-123456789-1234567890-1234567",
				UserAccountControlRaw: 8192,
				UserAccountControl:    []string{"SERVER_TRUST_ACCOUNT"},
			}},
	} {
		var got DomainController
		if err := got.Read(tt.args); err != nil {
			t.Error(err)
		}

		changelogs, err := diff.Diff(got, tt.want, diff.DisableStructValues())
		if err != nil {
			t.Error(err)
		}

		for _, changelog := range changelogs {
			t.Errorf(`(DomainController).Read(...) failed: Type: %s, Path: %v, Got: %v, Want: %v`, changelog.Type, changelog.Path, changelog.From, changelog.To)
		}
	}
}
