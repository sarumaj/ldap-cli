package objects

import (
	"fmt"
	"testing"
	"time"

	"github.com/r3labs/diff/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestUserRead(t *testing.T) {
	defer util.PatchForTimeNow().Unpatch()
	for _, tt := range []struct {
		name string
		args map[string]any
		want User
	}{
		{"test#1",
			map[string]any{
				"accountExpires":              "1234567891234567811",
				"badPwdCount":                 "3",
				"badPasswordTime":             "1234567891234567811",
				"c":                           "DE",
				"cn":                          "test",
				"co":                          "Germany",
				"company":                     "DB",
				"distinguishedName":           "CN=test,...DN:example,DN=com",
				"division":                    "test",
				"givenName":                   "John",
				"global-extensionAttribute11": "RegIntPri",
				"global-extensionAttribute22": "Berlin",
				"global-extensionAttribute26": "Black Street 1",
				"l":                           "Berlin",
				"mail":                        "milton@example.com",
				"memberOf":                    []string{"CN=parent,...DN:example,DN=com"},
				"msRadiusFramedIPAddress":     "2130706433",
				"name":                        "test",
				"objectCategory":              []string{"user"},
				"objectClass":                 []string{"user"},
				"objectGUID":                  "696768d7-5ca1-97a6-4835-cbab7c9e5b11",
				"objectSID":                   "S-1-5-21-1234567890-123456789-1234567890-1234567",
				"pwdLastSet":                  "1234567891234567811",
				"sn":                          "Milton",
				"SAMAccountName":              "test",
				"SAMAccountType":              fmt.Sprintf("%d", attributes.SAM_ACCOUNT_TYPE_USER_OBJECT),
				"userAccountControl":          fmt.Sprintf("%d", attributes.USER_ACCOUNT_CONTROL_LOCKOUT|attributes.USER_ACCOUNT_CONTROL_NORMAL_ACCOUNT),
				"userPrincipalName":           "test",
			},
			User{
				AccountExpires:             1234567891234567811,
				AccountExpiryDate:          time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC),
				BadPasswordCount:           3,
				BadPasswordTime:            time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC),
				BadPasswordTimeRaw:         1234567891234567811,
				CountryCode:                "DE",
				CountryName:                "Germany",
				CommonName:                 "test",
				Company:                    "DB",
				DistinguishedName:          "CN=test,...DN:example,DN=com",
				Division:                   "test",
				Enabled:                    true,
				GivenName:                  "John",
				GlobalExtensionAttribute11: "RegIntPri",
				GlobalExtensionAttribute22: "Berlin",
				GlobalExtensionAttribute26: "Black Street 1",
				Location:                   "Berlin",
				LockedOut:                  true,
				Mail:                       "milton@example.com",
				MemberOf:                   []string{"CN=parent,...DN:example,DN=com"},
				MsRadiusFramedIpAddressRaw: 2130706433,
				MsRadiusFramedIpAddress:    "127.0.0.1",
				Name:                       "test",
				ObjectCategory:             "user",
				ObjectClass:                []string{"user"},
				ObjectGUID:                 "696768d7-5ca1-97a6-4835-cbab7c9e5b11",
				PasswordLastSet:            time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC),
				PasswordLastSetRaw:         1234567891234567811,
				SamAccountName:             "test",
				SamAccountType:             []string{"NORMAL_USER_ACCOUNT", "USER_OBJECT"},
				SamAccountTypeRaw:          805306368,
				SID:                        "S-1-5-21-1234567890-123456789-1234567890-1234567",
				Surname:                    "Milton",
				UserAccountControlRaw:      528,
				UserAccountControl:         []string{"LOCKOUT", "NORMAL_ACCOUNT"},
				UserPrincipalName:          "test"}},
	} {
		var got User
		if err := got.Read(tt.args); err != nil {
			t.Error(err)
		}

		changelogs, err := diff.Diff(got, tt.want, diff.DisableStructValues())
		if err != nil {
			t.Error(err)
		}

		for _, changelog := range changelogs {
			t.Errorf(`(User).Read(...) failed: Type: %s, Path: %v, Got: %v, Want: %v`, changelog.Type, changelog.Path, changelog.From, changelog.To)
		}
	}
}
