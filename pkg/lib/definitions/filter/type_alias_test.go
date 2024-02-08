package filter

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	attributes "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
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
		{"test#1", args{AliasForEnabled(), ".....$ENABLED......."}, nil},
		{"test#2", args{AliasForDisabled(), ".....$DISABLED......."}, nil},
		{"test#3", args{AliasForGroup(), ".....$GROUP......."}, nil},
		{"test#4", args{AliasForUser(), ".....$USER.....$USERS.."}, nil},
		{"test#5", args{AliasForDc(), ".....$DC...$DC_FR...."}, nil},
		{"test#6", args{AliasForExpired(), ".....$EXPIRED......."}, nil},
		{"test#7", args{AliasForNotExpired(), ".....$NOT_EXPIRED......."}, nil},
		{"test#8", args{AliasForId(), ".....$ID(12345)......."}, []string{"$ID(12345)"}},
		{"test#9", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, []string{"$MEMBER_OF(CN=SuperUsers,...,DC=com)"}},
		{"test#10", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, []string{"$MEMBER_OF(CN=SuperUsers,...,DC=com;true)"}},
		{"test#11", args{AliasForAnd(), ".....$AND($AND($ID(12345)$ID(12346)))......."}, []string{"$AND($AND($ID(12345)$ID(12346)))", "$AND($ID(12345)$ID(12346))"}},
		{"test#12", args{AliasForNot(), ".....$NOT($NOT_EXPIRED)......."}, []string{"$NOT($NOT_EXPIRED)"}},
		{"test#13", args{AliasForOr(), ".....$OR($ID(12345)$ID(12346))......."}, []string{"$OR($ID(12345)$ID(12346))"}},
		{"test#14", args{AliasForBand(), ".....$BAND......."}, nil},
		{"test#15", args{AliasForBor(), ".....$BOR......."}, nil},
		{"test#16", args{AliasForRecursive(), ".....$RECURSIVE......."}, nil},
		{"test#17", args{AliasForData(), ".....$DATA......."}, nil},
		{"test#18", args{AliasForAttr(), ".....$ATTR()........"}, []string{"$ATTR()"}},
		{"test#19", args{AliasForAttr(), ".....$ATTR(dn)........"}, []string{"$ATTR(dn)"}},
		{"test#20", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com)"}},
		{"test#21", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com;=*)"}},
		{"test#22", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)"}},
		{"test#23", args{AliasForAttr(), ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, []string{"$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)"}},
		{"test#24", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)..$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"},
			[]string{"$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)", "$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)"}},
		{"test#25", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com)........"}, []string{"$EQ(dn;CN=SuperUser,...,DC=com)"}},
		{"test#26", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com;$BAND)........"}, []string{"$EQ(dn;CN=SuperUser,...,DC=com;$BAND)"}},
		{"test#27", args{AliasForGreaterThan(), ".....$GT(createdOn;0)........"}, []string{"$GT(createdOn;0)"}},
		{"test#28", args{AliasForGreaterThanOrEqual(), "....$GT(createdOn;0).$GTE(createdOn;0).$GT(createdOn;0)......."}, []string{"$GTE(createdOn;0)"}},
		{"test#29", args{AliasForLessThan(), ".....$LT(createdOn;0)........"}, []string{"$LT(createdOn;0)"}},
		{"test#30", args{AliasForLessThanOrEqual(), "....$LT(createdOn;0).$LTE(createdOn;0).$LT(createdOn;0)......."}, []string{"$LTE(createdOn;0)"}},
		{"test#31", args{AliasForContains(), ".....$CONTAINS(commonName;user)........"}, []string{"$CONTAINS(commonName;user)"}},
		{"test#32", args{AliasForStartsWith(), ".....$STARTS_WITH(commonName;user)........"}, []string{"$STARTS_WITH(commonName;user)"}},
		{"test#33", args{AliasForEndsWith(), ".....$ENDS_WITH(commonName;user)........"}, []string{"$ENDS_WITH(commonName;user)"}},
		{"test#34", args{AliasForNotExists(), ".....$NOT_EXISTS(commonName)........"}, []string{"$NOT_EXISTS(commonName)"}},
		{"test#35", args{AliasForExists(), ".....$EXISTS(commonName)........"}, []string{"$EXISTS(commonName)"}},
		{"test#36", args{AliasForLike(), ".....$LIKE(commonName;user)........"}, []string{"$LIKE(commonName;user)"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.alias.findMatches(tt.args.alias.findOccurrences([]byte(tt.args.raw)), []byte(tt.args.raw))
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
		{"test#1", args{AliasForEnabled(), ".....$ENABLED......."}, [][]int{{5, 13}}},
		{"test#2", args{AliasForDisabled(), ".....$DISABLED......."}, [][]int{{5, 14}}},
		{"test#3", args{AliasForGroup(), ".....$GROUP......."}, [][]int{{5, 11}}},
		{"test#4", args{AliasForUser(), ".....$USER.....$USERS.."}, [][]int{{5, 10}}},
		{"test#5", args{AliasForDc(), ".....$DC...$DC_FR...."}, [][]int{{5, 8}}},
		{"test#6", args{AliasForExpired(), ".....$EXPIRED......."}, [][]int{{5, 13}}},
		{"test#7", args{AliasForNotExpired(), ".....$NOT_EXPIRED......."}, [][]int{{5, 17}}},
		{"test#8", args{AliasForId(), ".....$ID(12345)......."}, [][]int{{5, 8}}},
		{"test#9", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, [][]int{{5, 15}}},
		{"test#10", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, [][]int{{5, 15}}},
		{"test#11", args{AliasForAnd(), ".....$AND($AND($ID(12345)$ID(12346)))......."}, [][]int{{5, 9}, {10, 14}}},
		{"test#12", args{AliasForNot(), ".....$NOT($NOT_EXPIRED)......."}, [][]int{{5, 9}}},
		{"test#13", args{AliasForOr(), ".....$OR($ID(12345)$ID(12346))......."}, [][]int{{5, 8}}},
		{"test#14", args{AliasForBand(), ".....$BAND......."}, [][]int{{5, 10}}},
		{"test#15", args{AliasForBor(), ".....$BOR......."}, [][]int{{5, 9}}},
		{"test#16", args{AliasForRecursive(), ".....$RECURSIVE......."}, [][]int{{5, 15}}},
		{"test#17", args{AliasForData(), ".....$DATA......."}, [][]int{{5, 10}}},
		{"test#18", args{AliasForNotExpired(), ".....$NOT.EXPIRED......."}, nil},
		{"test#19", args{AliasForNotExpired(), ".....$NOT.$EXPIRED......."}, nil},
		{"test#20", args{AliasForNot(), ".....$NOT.$EXPIRED......."}, [][]int{{5, 9}}},
		{"test#21", args{AliasForExpired(), ".....$NOT.$EXPIRED......."}, [][]int{{10, 18}}},
		{"test#22", args{AliasForExpired(), ".....$NOT.$EXPIRED......."}, [][]int{{10, 18}}},
		{"test#24", args{AliasForAttr(), ".....$ATTR()........"}, [][]int{{5, 10}}},
		{"test#25", args{AliasForAttr(), ".....$ATTR(dn)........"}, [][]int{{5, 10}}},
		{"test#26", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, [][]int{{5, 10}}},
		{"test#27", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, [][]int{{5, 10}}},
		{"test#28", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, [][]int{{5, 10}}},
		{"test#29", args{AliasForAttr(), ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, [][]int{{10, 15}}},
		{"test#30", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com)........"}, [][]int{{5, 8}}},
		{"test#31", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com;$BAND)........"}, [][]int{{5, 8}}},
		{"test#32", args{AliasForGreaterThan(), ".....$GT(createdOn;0)........"}, [][]int{{5, 8}}},
		{"test#33", args{AliasForGreaterThanOrEqual(), "....$GT(createdOn;0).$GTE(createdOn;0).$GT(createdOn;0)......."}, [][]int{{21, 25}}},
		{"test#34", args{AliasForLessThan(), ".....$LT(createdOn;0)........"}, [][]int{{5, 8}}},
		{"test#35", args{AliasForLessThanOrEqual(), "....$LT(createdOn;0).$LTE(createdOn;0).$LT(createdOn;0)......."}, [][]int{{21, 25}}},
		{"test#36", args{AliasForContains(), ".....$CONTAINS(commonName;user)........"}, [][]int{{5, 14}}},
		{"test#37", args{AliasForStartsWith(), ".....$STARTS_WITH(commonName;user)........"}, [][]int{{5, 17}}},
		{"test#38", args{AliasForEndsWith(), ".....$ENDS_WITH(commonName;user)........"}, [][]int{{5, 15}}},
		{"test#39", args{AliasForNotExists(), ".....$NOT_EXISTS(commonName)........"}, [][]int{{5, 16}}},
		{"test#40", args{AliasForExists(), ".....$EXISTS(commonName)........"}, [][]int{{5, 12}}},
		{"test#41", args{AliasForLike(), ".....$LIKE(commonName;user)........"}, [][]int{{5, 10}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.alias.findOccurrences([]byte(tt.args.raw)); !reflect.DeepEqual(got, tt.want) {
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
		{"test#1", args{AliasForEnabled(), ".....$ENABLED......."}, nil},
		{"test#2", args{AliasForDisabled(), ".....$DISABLED......."}, nil},
		{"test#3", args{AliasForGroup(), ".....$GROUP......."}, nil},
		{"test#4", args{AliasForUser(), ".....$USER.....$USERS.."}, nil},
		{"test#5", args{AliasForDc(), ".....$DC...$DC_FR...."}, nil},
		{"test#6", args{AliasForExpired(), ".....$EXPIRED......."}, nil},
		{"test#7", args{AliasForNotExpired(), ".....$NOT_EXPIRED......."}, nil},
		{"test#8", args{AliasForId(), ".....$ID(12345)......."}, []int{15}},
		{"test#9", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, []int{41}},
		{"test#10", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, []int{40, 46}},
		{"test#11", args{AliasForAnd(), ".....$AND($AND($ID(12345)$ID(12346)))......."}, []int{37}},
		{"test#12", args{AliasForNot(), ".....$NOT($NOT_EXPIRED)......."}, []int{23}},
		{"test#13", args{AliasForOr(), ".....$OR($ID(12345);$ID(12346))......."}, []int{19, 31}},
		{"test#14", args{AliasForBand(), ".....$BAND......."}, nil},
		{"test#15", args{AliasForBor(), ".....$BOR......."}, nil},
		{"test#16", args{AliasForRecursive(), ".....$RECURSIVE......."}, nil},
		{"test#17", args{AliasForData(), ".....$DATA......."}, nil},
		{"test#18", args{AliasForAnd(), ".....$AND($AND($ID(12345);$ID(12346)))......."}, []int{38}},
		{"test#19", args{AliasForAttr(), ".....$ATTR()........"}, []int{12}},
		{"test#20", args{AliasForAttr(), ".....$ATTR(dn)........"}, []int{14}},
		{"test#21", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, []int{13, 38}},
		{"test#22", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, []int{13, 37, 41}},
		{"test#23", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, []int{13, 37, 39, 46}},
		{"test#24", args{AliasForAttr(), ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, []int{52}},
		{"test#25", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)..$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, []int{13, 37, 39, 46}},
		{"test#26", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com)........"}, []int{11, 36}},
		{"test#27", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com;$BAND)........"}, []int{11, 35, 42}},
		{"test#28", args{AliasForGreaterThan(), ".....$GT(createdOn;0)........"}, []int{18, 21}},
		{"test#29", args{AliasForGreaterThanOrEqual(), "....$GT(createdOn;0).$GTE(createdOn;0).$GT(createdOn;0)......."}, []int{17, 20}},
		{"test#30", args{AliasForLessThan(), ".....$LT(createdOn;0)........"}, []int{18, 21}},
		{"test#31", args{AliasForLessThanOrEqual(), "....$LT(createdOn;0).$LTE(createdOn;0).$LT(createdOn;0)......."}, []int{17, 20}},
		{"test#32", args{AliasForContains(), ".....$CONTAINS(commonName;user)........"}, []int{25, 31}},
		{"test#33", args{AliasForStartsWith(), ".....$STARTS_WITH(commonName;user)........"}, []int{28, 34}},
		{"test#34", args{AliasForEndsWith(), ".....$ENDS_WITH(commonName;user)........"}, []int{26, 32}},
		{"test#35", args{AliasForNotExists(), ".....$NOT_EXISTS(commonName)........"}, []int{28}},
		{"test#36", args{AliasForExists(), ".....$EXISTS(commonName)........"}, []int{24}},
		{"test#37", args{AliasForLike(), ".....$LIKE(commonName;user)........"}, []int{21, 27}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.alias.findSplitPositions([]byte(tt.args.raw)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(Alias{ID: %q}).findSplitPositions() = %v, want %v", tt.args.alias.ID, got, tt.want)
			}
		})
	}
}

func TestAlias_replace(t *testing.T) {
	util.Now = func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) }
	t.Cleanup(func() { util.Now = time.Now })

	type args struct {
		alias Alias
		raw   string
	}

	for _, tt := range []struct {
		name string
		args args
		want string
	}{
		{"test#1", args{AliasForEnabled(), ".....$ENABLED......."}, ".....(!(UserAccountControl:1.2.840.113556.1.4.803:=2))......."},
		{"test#2", args{AliasForDisabled(), ".....$DISABLED......."}, ".....(UserAccountControl:1.2.840.113556.1.4.803:=2)......."},
		{"test#3", args{AliasForGroup(), ".....$GROUP......."}, ".....(|(ObjectClass=group)(ObjectClass=posixGroup))......."},
		{"test#4", args{AliasForUser(), ".....$USER.....$USERS.."}, ".....(|(ObjectClass=user)(ObjectClass=posixAccount)).....$USERS.."},
		{"test#5", args{AliasForDc(), ".....$DC...$DC_FR...."}, ".....(&(ObjectClass=computer)(UserAccountControl:1.2.840.113556.1.4.803:=8192))...$DC_FR...."},
		{"test#6", args{AliasForExpired(), ".....$EXPIRED......."}, ".....(&" +
			"(AccountExpires=>0)" +
			"(AccountExpires=<9223372036854775807)" +
			"(AccountExpires=<92233720368547758)" +
			"(AccountExpires=*)" +
			")......."},
		{"test#7", args{AliasForNotExpired(), ".....$NOT_EXPIRED......."}, ".....(!" +
			"(&" +
			"(AccountExpires=>0)" +
			"(AccountExpires=<9223372036854775807)" +
			"(AccountExpires=<92233720368547758)" +
			"(AccountExpires=*)" +
			")" +
			")......."},
		{"test#8", args{AliasForId(), ".....$ID(12345)......."}, ".....(|" +
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
		{"test#9", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com)......."}, ".....(MemberOf=CN=SuperUsers,...,DC=com)......."},
		{"test#10", args{AliasForMemberOf(), ".....$MEMBER_OF(CN=SuperUsers,...,DC=com;true)......."}, ".....(MemberOf:1.2.840.113556.1.4.1941:=CN=SuperUsers,...,DC=com)......."},
		{"test#11", args{AliasForAnd(), ".....$AND($AND($ID(12345)$ID(12346)))......."}, ".....(&(&$ID(12345)$ID(12346)))......."},
		{"test#12", args{AliasForNot(), ".....$NOT($NOT_EXPIRED)......."}, ".....(!$NOT_EXPIRED)......."},
		{"test#13", args{AliasForOr(), ".....$OR($ID(12345)$ID(12346))......."}, ".....(|$ID(12345)$ID(12346))......."},
		{"test#14", args{AliasForBand(), ".....$BAND......."}, "....." + string(attributes.LDAP_MATCHING_RULE_BIT_AND) + "......."},
		{"test#15", args{AliasForBor(), ".....$BOR......."}, "....." + string(attributes.LDAP_MATCHING_RULE_BIT_OR) + "......."},
		{"test#16", args{AliasForRecursive(), ".....$RECURSIVE......."}, "....." + string(attributes.LDAP_MATCHING_RULE_IN_CHAIN) + "......."},
		{"test#17", args{AliasForData(), ".....$DATA......."}, "....." + string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA) + "......."},
		{"test#18", args{AliasForNotExpired(), ".....$NOT.EXPIRED......."}, ".....$NOT.EXPIRED......."},
		{"test#19", args{AliasForNotExpired(), ".....$NOT.$EXPIRED......."}, ".....$NOT.$EXPIRED......."},
		{"test#20", args{AliasForNot(), ".....$NOT.$EXPIRED......."}, ".....$NOT.$EXPIRED......."},
		{"test#21", args{AliasForExpired(), ".....$NOT.$EXPIRED......."}, ".....$NOT.(&" +
			"(AccountExpires=>0)" +
			"(AccountExpires=<9223372036854775807)" +
			"(AccountExpires=<92233720368547758)" +
			"(AccountExpires=*)" +
			")......."},
		{"test#22", args{AliasForNotExpired(), ".....$NOT().............."}, ".....$NOT().............."},
		{"test#23", args{AliasForMemberOf(), ".....$MEMBER_OF()........"}, ".....(MemberOf=)........"},
		{"test#24", args{AliasForAttr(), ".....$ATTR()........"}, ".....(=)........"},
		{"test#25", args{AliasForAttr(), ".....$ATTR(dn)........"}, ".....(DistinguishedName=)........"},
		{"test#26", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com)........"}, ".....(DistinguishedName=CN=SuperUser,...,DC=com)........"},
		{"test#27", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=*)........"}, ".....(DistinguishedName=*CN=SuperUser,...,DC=com)........"},
		{"test#28", args{AliasForAttr(), ".....$ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST)........"}, ".....(DistinguishedName:$TEST:=CN=SuperUser,...,DC=com)........"},
		{"test#29", args{AliasForAttr(), ".....$NOT($ATTR(dn;CN=SuperUser,...,DC=com;=;$TEST))........"}, ".....$NOT((DistinguishedName:$TEST:=CN=SuperUser,...,DC=com))........"},
		{"test#30", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com)........"}, ".....(DistinguishedName=CN=SuperUser,...,DC=com)........"},
		{"test#31", args{AliasForEquals(), ".....$EQ(dn;CN=SuperUser,...,DC=com;$BAND)........"}, ".....(DistinguishedName:$BAND:=CN=SuperUser,...,DC=com)........"},
		{"test#32", args{AliasForGreaterThan(), ".....$GT(createdOn;0)........"}, ".....(CreatedOn=>0)........"},
		{"test#33", args{AliasForGreaterThanOrEqual(), "....$GT(createdOn;0).$GTE(createdOn;0).$GT(createdOn;0)......."},
			"....$GT(createdOn;0).(CreatedOn>=0).$GT(createdOn;0)......."},
		{"test#34", args{AliasForLessThan(), ".....$LT(createdOn;0)........"}, ".....(CreatedOn=<0)........"},
		{"test#35", args{AliasForLessThanOrEqual(), "....$LT(createdOn;0).$LTE(createdOn;0).$LT(createdOn;0)......."},
			"....$LT(createdOn;0).(CreatedOn<=0).$LT(createdOn;0)......."},
		{"test#36", args{AliasForContains(), ".....$CONTAINS(commonName;user)........"}, ".....(CN=*user*)........"},
		{"test#37", args{AliasForStartsWith(), ".....$STARTS_WITH(commonName;user)........"}, ".....(CN=user*)........"},
		{"test#38", args{AliasForEndsWith(), ".....$ENDS_WITH(commonName;user)........"}, ".....(CN=*user)........"},
		{"test#39", args{AliasForNotExists(), ".....$NOT_EXISTS(commonName)........"}, ".....(!(CN=*))........"},
		{"test#40", args{AliasForExists(), ".....$EXISTS(commonName)........"}, ".....(CN=*)........"},
		{"test#41", args{AliasForLike(), ".....$LIKE(commonName;user)........"}, ".....(CN~=user)........"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.alias.replace([]byte(tt.args.raw)); string(got) != tt.want {
				t.Errorf("(Alias{ID: %q}).findOccurrences() = %q, want %q", tt.args.alias.ID, got, tt.want)
			}
		})
	}
}

func TestAlias_splitParameters(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want []string
	}{
		{"test#1", "(12345)", []string{"12345"}},
		{"test#2", "(12345;67890)", []string{"12345", "67890"}},
		{"test#3", "(12345;67890;12345)", []string{"12345", "67890", "12345"}},
		{"test#4", "(12345;67890;12345;67890)", []string{"12345", "67890", "12345", "67890"}},
		{"test#5", "(12345;67890;12345; ;67890 ; 12345)", []string{"12345", "67890", "12345", "", "67890", "12345"}},
		{"test#6", "((12345;67890;12345;67890;12345);34567)", []string{"(12345;67890;12345;67890;12345)", "34567"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			alias := Alias{}
			if got := alias.splitParameters([]byte(tt.args), alias.findSplitPositions([]byte(tt.args))); !reflect.DeepEqual(got, tt.want) {
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
	util.Now = func() time.Time { return time.Date(2023, 9, 12, 16, 27, 13, 0, time.UTC) }
	t.Cleanup(func() { util.Now = time.Now })

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
		{"test#8", "$OR($ID(12345);$ID(12346))",
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
		{"test#18", "$AND($GT(expiresOn;0);$LTE(expiresOn;9999);$EXISTS(expiresOn))", "(&(ExpiresOn=>0)(ExpiresOn<=9999)(ExpiresOn=*))"},
		{"test#19", "$AND($CONTAINS(cn;user);$STARTS_WITH(cn;prefix);$ENDS_WITH(cn;suffix))", "(&(CN=*user*)(CN=prefix*)(CN=*suffix))"},
		{"test#20", "$AND($NOT_EXISTS(cn);$EXISTS(cn))", "(&(!(CN=*))(CN=*))"},
		{"test#21", "$AND($LIKE(cn;user);$EQ(cn;user);$GT(cn;user);$GTE(cn;user);$LT(cn;user);$LTE(cn;user))",
			"(&(CN~=user)(CN=user)(CN=>user)(CN=user)(CN=<user)(CN=user))"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceAliases(tt.args); got != tt.want {
				t.Errorf("ReplaceAliases() = %v, want %v", got, tt.want)
			}
		})
	}
}
