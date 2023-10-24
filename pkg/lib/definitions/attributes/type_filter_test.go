package attributes

import (
	"testing"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestFilter(t *testing.T) {
	defer util.PatchForTimeNow().Unpatch()
	for _, tt := range []struct {
		name string
		args Filter
		want string
	}{
		{"test#1", And(
			Filter{"objectClass", "User", ""},
			HasNotExpired(true),
			Or(
				Filter{"sAMAccountName", "uid12345", ""},
				Filter{"sAMAccountName", "uid54321", ""},
			),
		), `(&` +
			(`` +
				`(objectClass=User)` +
				`(&` +
				(`` +
					`(|` +
					(`` +
						`(accountExpires=0)` +
						`(accountExpires=9223372036854775807)` + `(accountExpires>=92233720368547758)`) +
					`)` +
					`(accountExpires=*)`) +
				`)` +
				`(|` +
				(`` +
					`(sAMAccountName=uid12345)` +
					`(sAMAccountName=uid54321)`) +
				`)`) +
			`)`},
		{"test#2", And(
			Filter{"objectClass", "User", ""},
			Filter{"sAMAccountName", "uid12345", ""},
			Not(Filter{"userAccountControl", "2", LDAP_MATCHING_RULE_BIT_AND}),
			Not(Filter{
				"memberOf",
				"CN=SuperUsers,...,DC=com",
				LDAP_MATCHING_RULE_IN_CHAIN,
			})), `(&` +
			(`` +
				`(objectClass=User)` +
				`(sAMAccountName=uid12345)` +
				`(!(userAccountControl:1.2.840.113556.1.4.803:=2))` +
				`(!(memberOf:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com))`) +
			`)`},
		{"test#3", And(
			Filter{"displayName", EscapeFilter("id@dom"), ""},
			HasNotExpired(false),
			Not(IsDomainController()),
			Not(IsGroup()),
			IsUser(),
			IsEnabled()), `(&` +
			(`` +
				`(displayName=id@dom)` +
				`(|` +
				`(|` +
				(`` +
					(`(accountExpires=0)` +
						`(accountExpires=9223372036854775807)` +
						`(accountExpires>=92233720368547758)`) +
					`)` +
					`(!(accountExpires=*))`) +
				`)` +
				`(!(&` +
				(`` + (`` +
					`(objectClass=computer)` +
					`(userAccountControl:1.2.840.113556.1.4.803:=8192)`)) +
				`))` +
				`(!(objectClass=group))` +
				`(objectClass=user)` +
				`(!(userAccountControl:1.2.840.113556.1.4.803:=2))`) +
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
