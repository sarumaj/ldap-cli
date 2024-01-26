package filter

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	attributes "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/attributes"
)

func TestAlias_findMatches(t *testing.T) {
	type args struct {
		alias Alias
		raw   string
	}

	for _, tt := range []struct {
		name string
		args args
		want []string
	}{
		{"test#1", args{enabled, ".....$ENABLED......."}, nil},
		{"test#2", args{disabled, ".....$DISABLED......."}, nil},
		{"test#3", args{group, ".....$GROUP......."}, nil},
		{"test#4", args{user, ".....$USER.....$USERS.."}, nil},
		{"test#5", args{dc, ".....$DC...$DC_FR...."}, nil},
		{"test#6", args{expired, ".....$EXPIRED......."}, nil},
		{"test#7", args{not_expired, ".....$NOT_EXPIRED......."}, nil},
		{"test#8", args{id, ".....$ID(12345)......."}, []string{"$ID(12345)"}},
		{"test#9", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, []string{"$MEMBER_OF(CN=SuperUsers,...,DC=com)"}},
		{"test#10", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, []string{"$MEMBER_OF(CN=SuperUsers,...,DC=com;true)"}},
		{"test#11", args{and, ".....$AND($AND($ID(12345)$ID(12346)))......."}, []string{"$AND($AND($ID(12345)$ID(12346)))", "$AND($ID(12345)$ID(12346))"}},
		{"test#12", args{not, ".....$NOT($NOT_EXPIRED)......."}, []string{"$NOT($NOT_EXPIRED)"}},
		{"test#13", args{or, ".....$OR($ID(12345)$ID(12346))......."}, []string{"$OR($ID(12345)$ID(12346))"}},
		{"test#14", args{band, ".....$BAND......."}, nil},
		{"test#15", args{bor, ".....$BOR......."}, nil},
		{"test#16", args{recursive, ".....$RECURSIVE......."}, nil},
		{"test#17", args{data, ".....$DATA......."}, nil},
		{"test#18", args{attr, ".....$ATTR()........"}, []string{"$ATTR()"}},
		{"test#19", args{attr, ".....$ATTR(dn)........"}, []string{"$ATTR(dn)"}},
		{"test#20", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com)"}},
		{"test#21", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com;=*)"}},
		{"test#22", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)"}},
		{"test#23", args{attr, ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)"}},
		{"test#24", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)..$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"},
			[]string{"$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)", "$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.alias.findMatches(tt.args.alias.findOccurences([]byte(tt.args.raw)), []byte(tt.args.raw))
			if gotS, wantS := `[`+string(bytes.Join(got, []byte{' '}))+`]`, fmt.Sprint(tt.want); gotS != wantS {
				t.Errorf("(Alias{ID: %q}).findMatches() = %q, want %q", tt.args.alias.ID, gotS, wantS)
			}
		})
	}
}

func TestAlias_findOccurrences(t *testing.T) {
	type args struct {
		alias Alias
		raw   string
	}

	for _, tt := range []struct {
		name string
		args args
		want [][]int
	}{
		{"test#1", args{enabled, ".....$ENABLED......."}, [][]int{{5, 13}}},
		{"test#2", args{disabled, ".....$DISABLED......."}, [][]int{{5, 14}}},
		{"test#3", args{group, ".....$GROUP......."}, [][]int{{5, 11}}},
		{"test#4", args{user, ".....$USER.....$USERS.."}, [][]int{{5, 10}}},
		{"test#5", args{dc, ".....$DC...$DC_FR...."}, [][]int{{5, 8}}},
		{"test#6", args{expired, ".....$EXPIRED......."}, [][]int{{5, 13}}},
		{"test#7", args{not_expired, ".....$NOT_EXPIRED......."}, [][]int{{5, 17}}},
		{"test#8", args{id, ".....$ID(12345)......."}, [][]int{{5, 8}}},
		{"test#9", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, [][]int{{5, 15}}},
		{"test#10", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, [][]int{{5, 15}}},
		{"test#11", args{and, ".....$AND($AND($ID(12345)$ID(12346)))......."}, [][]int{{5, 9}, {10, 14}}},
		{"test#12", args{not, ".....$NOT($NOT_EXPIRED)......."}, [][]int{{5, 9}}},
		{"test#13", args{or, ".....$OR($ID(12345)$ID(12346))......."}, [][]int{{5, 8}}},
		{"test#14", args{band, ".....$BAND......."}, [][]int{{5, 10}}},
		{"test#15", args{bor, ".....$BOR......."}, [][]int{{5, 9}}},
		{"test#16", args{recursive, ".....$RECURSIVE......."}, [][]int{{5, 15}}},
		{"test#17", args{data, ".....$DATA......."}, [][]int{{5, 10}}},
		{"test#18", args{not_expired, ".....$NOT.EXPIRED......."}, nil},
		{"test#19", args{not_expired, ".....$NOT.$EXPIRED......."}, nil},
		{"test#20", args{not, ".....$NOT.$EXPIRED......."}, [][]int{{5, 9}}},
		{"test#21", args{expired, ".....$NOT.$EXPIRED......."}, [][]int{{10, 18}}},
		{"test#22", args{expired, ".....$NOT.$EXPIRED......."}, [][]int{{10, 18}}},
		{"test#24", args{attr, ".....$ATTR()........"}, [][]int{{5, 10}}},
		{"test#25", args{attr, ".....$ATTR(dn)........"}, [][]int{{5, 10}}},
		{"test#26", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, [][]int{{5, 10}}},
		{"test#27", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, [][]int{{5, 10}}},
		{"test#28", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, [][]int{{5, 10}}},
		{"test#29", args{attr, ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, [][]int{{10, 15}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.alias.findOccurences([]byte(tt.args.raw)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(Alias{ID: %q}).findOccurrences() = %v, want %v", tt.args.alias.ID, got, tt.want)
			}
		})
	}
}

func TestAlias_findSplitPositions(t *testing.T) {
	type args struct {
		alias Alias
		raw   string
	}

	for _, tt := range []struct {
		name string
		args args
		want []int
	}{
		{"test#1", args{enabled, ".....$ENABLED......."}, nil},
		{"test#2", args{disabled, ".....$DISABLED......."}, nil},
		{"test#3", args{group, ".....$GROUP......."}, nil},
		{"test#4", args{user, ".....$USER.....$USERS.."}, nil},
		{"test#5", args{dc, ".....$DC...$DC_FR...."}, nil},
		{"test#6", args{expired, ".....$EXPIRED......."}, nil},
		{"test#7", args{not_expired, ".....$NOT_EXPIRED......."}, nil},
		{"test#8", args{id, ".....$ID(12345)......."}, []int{15}},
		{"test#9", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, []int{41}},
		{"test#10", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, []int{40, 46}},
		{"test#11", args{and, ".....$AND($AND($ID(12345)$ID(12346)))......."}, []int{37}},
		{"test#12", args{not, ".....$NOT($NOT_EXPIRED)......."}, []int{23}},
		{"test#13", args{or, ".....$OR($ID(12345);$ID(12346))......."}, []int{19, 31}},
		{"test#14", args{band, ".....$BAND......."}, nil},
		{"test#15", args{bor, ".....$BOR......."}, nil},
		{"test#16", args{recursive, ".....$RECURSIVE......."}, nil},
		{"test#17", args{data, ".....$DATA......."}, nil},
		{"test#18", args{and, ".....$AND($AND($ID(12345);$ID(12346)))......."}, []int{38}},
		{"test#19", args{attr, ".....$ATTR()........"}, []int{12}},
		{"test#20", args{attr, ".....$ATTR(dn)........"}, []int{14}},
		{"test#21", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, []int{13, 38}},
		{"test#22", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, []int{13, 37, 41}},
		{"test#23", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, []int{13, 37, 39, 46}},
		{"test#24", args{attr, ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, []int{52}},
		{"test#25", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)..$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, []int{13, 37, 39, 46}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.alias.findSplitPositions([]byte(tt.args.raw)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(Alias{ID: %q}).findSplitPositions() = %v, want %v", tt.args.alias.ID, got, tt.want)
			}
		})
	}
}

func TestAlias_replace(t *testing.T) {
	defer monkey.Patch(
		time.Now,
		func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) },
	).Unpatch()

	type args struct {
		alias Alias
		raw   string
	}

	for _, tt := range []struct {
		name string
		args args
		want string
	}{
		{"test#1", args{enabled, ".....$ENABLED......."}, ".....(!(UserAccountControl:1.2.840.113556.1.4.803:=2))......."},
		{"test#2", args{disabled, ".....$DISABLED......."}, ".....(UserAccountControl:1.2.840.113556.1.4.803:=2)......."},
		{"test#3", args{group, ".....$GROUP......."}, ".....(|(ObjectClass=group)(ObjectClass=posixGroup))......."},
		{"test#4", args{user, ".....$USER.....$USERS.."}, ".....(|(ObjectClass=user)(ObjectClass=posixAccount)).....$USERS.."},
		{"test#5", args{dc, ".....$DC...$DC_FR...."}, ".....(&(ObjectClass=computer)(UserAccountControl:1.2.840.113556.1.4.803:=8192))...$DC_FR...."},
		{"test#6", args{expired, ".....$EXPIRED......."}, ".....(&" +
			"(AccountExpires=>0)" +
			"(AccountExpires=<9223372036854775807)" +
			"(AccountExpires=<92233720368547758)" +
			"(AccountExpires=*)" +
			")......."},
		{"test#7", args{not_expired, ".....$NOT_EXPIRED......."}, ".....(!" +
			"(&" +
			"(AccountExpires=>0)" +
			"(AccountExpires=<9223372036854775807)" +
			"(AccountExpires=<92233720368547758)" +
			"(AccountExpires=*)" +
			")" +
			")......."},
		{"test#8", args{id, ".....$ID(12345)......."}, ".....(|" +
			"(CN=12345)" +
			"(DisplayName=12345)" +
			"(|" +
			"(DistinguishedName=12345)" +
			"(DN=12345)" +
			")" +
			"(Name=12345)" +
			"(SAMAccountName=12345)" +
			"(UserPrincipalName=12345)" +
			"(ObjectGuid=12345)" +
			")......."},
		{"test#9", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, ".....(MemberOf=CN=SuperUsers,...,DC=com)......."},
		{"test#10", args{member_of, ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, ".....(MemberOf:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com)......."},
		{"test#11", args{and, ".....$AND($AND($ID(12345)$ID(12346)))......."}, ".....(&(&$ID(12345)$ID(12346)))......."},
		{"test#12", args{not, ".....$NOT($NOT_EXPIRED)......."}, ".....(!$NOT_EXPIRED)......."},
		{"test#13", args{or, ".....$OR($ID(12345)$ID(12346))......."}, ".....(|$ID(12345)$ID(12346))......."},
		{"test#14", args{band, ".....$BAND......."}, "....." + string(attributes.LDAP_MATCHING_RULE_BIT_AND) + "......."},
		{"test#15", args{bor, ".....$BOR......."}, "....." + string(attributes.LDAP_MATCHING_RULE_BIT_OR) + "......."},
		{"test#16", args{recursive, ".....$RECURSIVE......."}, "....." + string(attributes.LDAP_MATCHING_RULE_IN_CHAIN) + "......."},
		{"test#17", args{data, ".....$DATA......."}, "....." + string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA) + "......."},
		{"test#18", args{not_expired, ".....$NOT.EXPIRED......."}, ".....$NOT.EXPIRED......."},
		{"test#19", args{not_expired, ".....$NOT.$EXPIRED......."}, ".....$NOT.$EXPIRED......."},
		{"test#20", args{not, ".....$NOT.$EXPIRED......."}, ".....$NOT.$EXPIRED......."},
		{"test#21", args{expired, ".....$NOT.$EXPIRED......."}, ".....$NOT.(&" +
			"(AccountExpires=>0)" +
			"(AccountExpires=<9223372036854775807)" +
			"(AccountExpires=<92233720368547758)" +
			"(AccountExpires=*)" +
			")......."},
		{"test#22", args{not_expired, ".....$NOT().............."}, ".....$NOT().............."},
		{"test#23", args{member_of, ".....$MEMBER_OF()........"}, ".....(MemberOf=)........"},
		{"test#24", args{attr, ".....$ATTR()........"}, ".....(=)........"},
		{"test#25", args{attr, ".....$ATTR(dn)........"}, ".....(DistinguishedName=)........"},
		{"test#26", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, ".....(DistinguishedName=CN=SuperUser,...,DC=com)........"},
		{"test#27", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, ".....(DistinguishedName=*CN=SuperUser,...,DC=com)........"},
		{"test#28", args{attr, ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, ".....(DistinguishedName:$TEST:=CN=SuperUser,...,DC=com)........"},
		{"test#29", args{attr, ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, ".....$NOT((DistinguishedName:$TEST:=CN=SuperUser,...,DC=com))........"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.alias.replace([]byte(tt.args.raw)); string(got) != tt.want {
				t.Errorf("(Alias{ID: %q}).findOccurrences() = %q, want %q", tt.args.alias.ID, got, tt.want)
			}
		})
	}
}

func TestAlias_splitParameters(t *testing.T) {
	type args struct {
		alias  Alias
		params string
	}

	for _, tt := range []struct {
		name string
		args args
		want []string
	}{
		{"test#1", args{id, "(12345)"}, []string{"12345"}},
		{"test#2", args{id, "(12345;67890)"}, []string{"12345", "67890"}},
		{"test#3", args{id, "(12345;67890;12345)"}, []string{"12345", "67890", "12345"}},
		{"test#4", args{id, "(12345;67890;12345;67890)"}, []string{"12345", "67890", "12345", "67890"}},
		{"test#5", args{id, "(12345;67890;12345; ;67890 ; 12345)"}, []string{"12345", "67890", "12345", "", "67890", "12345"}},
		{"test#6", args{id, "((12345;67890;12345;67890;12345);34567)"}, []string{"(12345;67890;12345;67890;12345)", "34567"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.alias.splitParameters([]byte(tt.args.params), tt.args.alias.findSplitPositions([]byte(tt.args.params))); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListAliases(t *testing.T) {
	got := ListAliases()
	if len(got) == 0 {
		t.Errorf("ListAliases() = %v, want %v", got, "not empty")
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
		{"test#8", "$OR($ID(12345)$ID(12346))",
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
					`(UserPrincipalName=12345)` +
					`(ObjectGuid=12345)`) +
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
					`(UserPrincipalName=12346)` +
					`(ObjectGuid=12346)`) +
				`)`) +
				`)`},
		{"test#9", "$BAND", string(attributes.LDAP_MATCHING_RULE_BIT_AND)},
		{"test#10", "$BOR", string(attributes.LDAP_MATCHING_RULE_BIT_OR)},
		{"test#11", "$RECURSIVE", string(attributes.LDAP_MATCHING_RULE_IN_CHAIN)},
		{"test#12", "$DATA", string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA)},
		{"test#13", "$MEMBER_OF(CN=SuperUsers,...,DC=com)", "(MemberOf=CN=SuperUsers,...,DC=com)"},
		{"test#14", "$AND($MEMBER_OF(CN=SuperUsers,...,DC=com);$MEMBER_OF(CN=LocalUsers,...,DC=com))",
			(`(&` +
				`(MemberOf=CN=SuperUsers,...,DC=com)` +
				`(MemberOf=CN=LocalUsers,...,DC=com)`) +
				`)`},
		{"test#15", "$MEMBER_OF(CN=SuperUsers,...,DC=com;true)", "(MemberOf:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com)"},
		{"test#16", "$AND($AND($MEMBER_OF(CN=SuperUsers,...,DC=com;true);$NOT($NOT_EXPIRED)))",
			(`(&` +
				(`(&` +
					`(MemberOf:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com)` +
					(`(!` +
						(`(!` +
							(`(&` +
								`(AccountExpires=>0)` +
								`(AccountExpires=<9223372036854775807)` +
								`(AccountExpires=<92233720368547758)` +
								`(AccountExpires=*)`) +
							`)`) +
						`)`) +
					`)`) +
				`)`) +
				`)`},
		{"test#17", "$NOT($ATTR(test;value))", "(!(Test=value))"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceAliases(tt.args); got != tt.want {
				t.Errorf("ReplaceAliases() = %v, want %v", got, tt.want)
			}
		})
	}
}
