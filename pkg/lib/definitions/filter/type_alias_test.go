package filter

import (
	"testing"
	"time"

	"bou.ke/monkey"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

func TestListAliases(t *testing.T) {
	got, want := ListAliases(), make([]Alias, len(aliases))
	_ = copy(want, aliases)

	if len(got) != len(want) {
		t.Errorf("ListAliases() failed: got: %v, want: %v", got, want)
	}
}

func TestReplaceAliases(t *testing.T) {
	defer monkey.Patch(
		time.Now,
		func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) },
	).Unpatch()

	for _, tt := range []struct {
		name string
		args string
		want string
	}{
		{"test#1", "$ENABLED", "(!(UserAccountControl:1.2.840.113556.1.4.803:=2))"},
		{"test#2", "$DISABLED", "(UserAccountControl:1.2.840.113556.1.4.803:=2)"},
		{"test#3", "$GROUP", "(|(ObjectClass=group)(ObjectClass=posixGroup))"},
		{"test#4", "$USER", "(|(ObjectClass=user)(ObjectClass=posixAccount))"},
		{"test#5", "$DC", "(&(ObjectClass=computer)(UserAccountControl:1.2.840.113556.1.4.803:=8192))"},
		{"test#6", "$EXPIRED",
			(`(&` +
				`(AccountExpires=>0)` +
				`(AccountExpires=<9223372036854775807)` +
				`(AccountExpires=<92233720368547758)` +
				`(AccountExpires=*)`) +
				`)`},
		{"test#7", "$NOT_EXPIRED",
			(`(!` +
				(`(&` +
					`(AccountExpires=>0)` +
					`(AccountExpires=<9223372036854775807)` +
					`(AccountExpires=<92233720368547758)` +
					`(AccountExpires=*)`) +
				`)`) +
				`)`},
		{"test#8", "(|$ID(12345)$ID(12346))",
			(`(|` +
				(`(|` +
					`(CN=12345)` +
					`(DisplayName=12345)` +
					(`(|` +
						`(DistinguishedName=12345)` +
						`(DN=12345)`) +
					`)` +
					`(Name=12345)` +
					`(SAMAccountName=12345)` +
					`(UserPrincipalName=12345)`) +
				`)` +
				(`(|` +
					`(CN=12346)` +
					`(DisplayName=12346)` +
					(`(|` +
						`(DistinguishedName=12346)` +
						`(DN=12346)`) +
					`)` +
					`(Name=12346)` +
					`(SAMAccountName=12346)` +
					`(UserPrincipalName=12346)`) +
				`)`) +
				`)`},
		{"test#9", "$AND", string(attributes.LDAP_MATCHING_RULE_BIT_AND)},
		{"test#10", "$OR", string(attributes.LDAP_MATCHING_RULE_BIT_OR)},
		{"test#11", "$RECURSIVE", string(attributes.LDAP_MATCHING_RULE_IN_CHAIN)},
		{"test#12", "$DATA", string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA)},
		{"test#13", "$MEMBER_OF(CN=SuperUsers,...,DC=com)", "(MemberOf=CN=SuperUsers,...,DC=com)"},
		{"test#14", "(&$MEMBER_OF(CN=SuperUsers,...,DC=com)$MEMBER_OF(CN=LocalUsers,...,DC=com))",
			(`(&` +
				`(MemberOf=CN=SuperUsers,...,DC=com)` +
				`(MemberOf=CN=LocalUsers,...,DC=com)`) +
				`)`},
		{"test#15", "$MEMBER_OF_RECURSIVE(CN=SuperUsers,...,DC=com)", "(MemberOf:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com)"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceAliases(tt.args); got != tt.want {
				t.Errorf("ReplaceAliases() = %v, want %v", got, tt.want)
			}
		})
	}
}
