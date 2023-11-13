package filter

import (
	"testing"
	"time"

	"bou.ke/monkey"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

func TestFilter(t *testing.T) {
	defer monkey.Patch(
		time.Now,
		func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) },
	).Unpatch()

	for _, tt := range []struct {
		name string
		args Filter
		want string
	}{
		{"test#1", And(
			Filter{attributes.ObjectClass(), "User", ""},
			HasNotExpired(true),
			Or(
				Filter{attributes.SamAccountName(), "uid12345", ""},
				Filter{attributes.SamAccountName(), "uid54321", ""},
				Filter{attributes.Raw("boolean", "", attributes.TypeBool), "true", ""},
			),
		), `(&` +
			(`` +
				`(ObjectClass=User)` +
				`(&` +
				(`` +
					`(|` +
					(`` +
						`(AccountExpires=0)` +
						`(AccountExpires=9223372036854775807)` + `(AccountExpires>=92233720368547758)`) +
					`)` +
					`(AccountExpires=*)`) +
				`)` +
				`(|` +
				(`` +
					`(SAMAccountName=uid12345)` +
					`(SAMAccountName=uid54321)` +
					`(Boolean=TRUE)`) +
				`)`) +
			`)`},
		{"test#2", And(
			Filter{attributes.ObjectClass(), "User", ""},
			Filter{attributes.SamAccountName(), "uid12345", ""},
			Not(Filter{attributes.UserAccountControl(), "2", attributes.LDAP_MATCHING_RULE_BIT_AND}),
			Not(Filter{
				attributes.MemberOf(),
				"CN=SuperUsers,...,DC=com",
				attributes.LDAP_MATCHING_RULE_IN_CHAIN,
			})), `(&` +
			(`` +
				`(ObjectClass=User)` +
				`(SAMAccountName=uid12345)` +
				`(!(UserAccountControl:1.2.840.113556.1.4.803:=2))` +
				`(!(MemberOf:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com))`) +
			`)`},
		{"test#3", And(
			Filter{attributes.DisplayName(), EscapeFilter("id@dom"), ""},
			HasNotExpired(false).ExpandAlias(),
			Or(Not(IsDomainController())),
			And(Not(IsGroup())),
			IsUser(),
			IsEnabled()), `(&` +
			(`` +
				`(DisplayName=id@dom)` +
				`(|` +
				`(|` +
				(`` +
					(`` +
						`(AccountExpires=0)` +
						`(AccountExpires=9223372036854775807)` +
						`(AccountExpires>=92233720368547758)`) +
					`)` +
					`(!(AccountExpires=*))`) +
				`)` +
				`(!(&` +
				(`` + (`` +
					`(ObjectClass=computer)` +
					`(UserAccountControl:1.2.840.113556.1.4.803:=8192)`)) +
				`))` +
				(`(!(|` +
					`(ObjectClass=group)` +
					`(ObjectClass=posixGroup)`) +
				`))` +
				(`(|` +
					`(ObjectClass=user)` +
					`(ObjectClass=posixAccount)`) +
				`)` +
				`(!(UserAccountControl:1.2.840.113556.1.4.803:=2))`) +
			`)`},
		{"test#4", And(
			ByID("test"),
			MemberOf("test#1", true),
			MemberOf("test#2", false)), `(&` +
			(`` +
				`(|` +
				(`` +
					`(CN=test)` +
					`(DisplayName=test)` +
					`(|` +
					(`` +
						`(DistinguishedName=test)` +
						`(DN=test)`) +
					`)` +
					`(Name=test)` +
					`(SAMAccountName=test)` +
					`(UserPrincipalName=test)`) +
				`)` +
				`(MemberOf:1.2.840.113556.1.4.1941:=test#1)` +
				`(MemberOf=test#2)`) +
			`)`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.String()
			if got != tt.want {
				t.Errorf(`(Filter).String() failed: got %q, want: %q`, got, tt.want)
			}
		})
	}
}
