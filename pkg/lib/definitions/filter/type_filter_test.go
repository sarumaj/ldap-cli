package filter

import (
	"testing"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestFilter(t *testing.T) {
	defer libutil.PatchForTimeNow().Unpatch()

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
				`(objectclass=User)` +
				`(&` +
				(`` +
					`(|` +
					(`` +
						`(accountexpires=0)` +
						`(accountexpires=9223372036854775807)` + `(accountexpires>=92233720368547758)`) +
					`)` +
					`(accountexpires=*)`) +
				`)` +
				`(|` +
				(`` +
					`(samaccountname=uid12345)` +
					`(samaccountname=uid54321)` +
					`(boolean=TRUE)`) +
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
				`(objectclass=User)` +
				`(samaccountname=uid12345)` +
				`(!(useraccountcontrol:1.2.840.113556.1.4.803:=2))` +
				`(!(memberof:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com))`) +
			`)`},
		{"test#3", And(
			Filter{attributes.DisplayName(), EscapeFilter("id@dom"), ""},
			HasNotExpired(false),
			Not(IsDomainController()),
			Not(IsGroup()),
			IsUser(),
			IsEnabled()), `(&` +
			(`` +
				`(displayname=id@dom)` +
				`(|` +
				`(|` +
				(`` +
					(`(accountexpires=0)` +
						`(accountexpires=9223372036854775807)` +
						`(accountexpires>=92233720368547758)`) +
					`)` +
					`(!(accountexpires=*))`) +
				`)` +
				`(!(&` +
				(`` + (`` +
					`(objectclass=computer)` +
					`(useraccountcontrol:1.2.840.113556.1.4.803:=8192)`)) +
				`))` +
				`(!(objectclass=group))` +
				`(objectclass=user)` +
				`(!(useraccountcontrol:1.2.840.113556.1.4.803:=2))`) +
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
