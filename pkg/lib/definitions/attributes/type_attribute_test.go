package attributes

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/r3labs/diff/v3"
)

func TestAttributeParse(t *testing.T) {
	type args struct {
		a Attribute
		v []string
	}
	for _, tt := range []struct {
		name string
		args args
		want AttributeMap
	}{
		{"test#1",
			args{attributeAccountExpires, []string{"128271382742968750"}},
			AttributeMap{AttributeAccountExpires(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#2",
			args{attributeBadPasswordTime, []string{"128271382742968750"}},
			AttributeMap{AttributeBadPasswordTime(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#3",
			args{attributeBadPasswordCount, []string{"3"}},
			AttributeMap{AttributeBadPasswordCount(): int64(3)}},
		{"test#4",
			args{attributeCommonName, []string{"test"}},
			AttributeMap{AttributeCommonName(): "test"}},
		{"test#5",
			args{attributeCountryCode, []string{"test"}},
			AttributeMap{AttributeCountryCode(): "test"}},
		{"test#6",
			args{attributeCountryName, []string{"test"}},
			AttributeMap{AttributeCountryName(): "test"}},
		{"test#7",
			args{attributeCountryNumber, []string{"103"}},
			AttributeMap{AttributeCountryNumber(): int64(103)}},
		{"test#8",
			args{attributeCompany, []string{"test"}},
			AttributeMap{AttributeCompany(): "test"}},
		{"test#9",
			args{attributeDepartment, []string{"test"}},
			AttributeMap{AttributeDepartment(): "test"}},
		{"test#10",
			args{attributeDepartmentNumber, []string{"134"}},
			AttributeMap{AttributeDepartmentNumber(): "134"}},
		{"test#11",
			args{attributeDescription, []string{"test"}},
			AttributeMap{AttributeDescription(): "test"}},
		{"test#12",
			args{attributeDisplayName, []string{"test"}},
			AttributeMap{AttributeDisplayName(): "test"}},
		{"test#13",
			args{attributeDistinguishedName, []string{"test"}},
			AttributeMap{AttributeDistinguishedName(): "test"}},
		{"test#14",
			args{attributeDivision, []string{"test"}},
			AttributeMap{AttributeDivision(): "test"}},
		{"test#15",
			args{attributeDNSHostname, []string{"test"}},
			AttributeMap{AttributeDNSHostname(): "test"}},
		{"test#16",
			args{attributeEmployeeID, []string{"test"}},
			AttributeMap{AttributeEmployeeID(): "test"}},
		{"test#17",
			args{attributeGivenName, []string{"test"}},
			AttributeMap{AttributeGivenName(): "test"}},
		{"test#18",
			args{attributeGlobalExtension1, []string{"test"}},
			AttributeMap{AttributeGlobalExtension1(): "test"}},
		{"test#19",
			args{attributeGlobalExtension2, []string{"test"}},
			AttributeMap{AttributeGlobalExtension2(): "test"}},
		{"test#20",
			args{attributeGlobalExtension3, []string{"test"}},
			AttributeMap{AttributeGlobalExtension3(): "test"}},
		{"test#21",
			args{attributeGlobalExtension4, []string{"test"}},
			AttributeMap{AttributeGlobalExtension4(): "test"}},
		{"test#22",
			args{attributeGlobalExtension5, []string{"test"}},
			AttributeMap{AttributeGlobalExtension5(): "test"}},
		{"test#23",
			args{attributeGlobalExtension6, []string{"test"}},
			AttributeMap{AttributeGlobalExtension6(): "test"}},
		{"test#24",
			args{attributeGlobalExtension7, []string{"test"}},
			AttributeMap{AttributeGlobalExtension7(): "test"}},
		{"test#25",
			args{attributeGlobalExtension8, []string{"test"}},
			AttributeMap{AttributeGlobalExtension8(): "test"}},
		{"test#26",
			args{attributeGlobalExtension9, []string{"test"}},
			AttributeMap{AttributeGlobalExtension9(): "test"}},
		{"test#27",
			args{attributeGlobalExtension10, []string{"test"}},
			AttributeMap{AttributeGlobalExtension10(): "test"}},
		{"test#28",
			args{attributeGlobalExtension11, []string{"test"}},
			AttributeMap{AttributeGlobalExtension11(): "test"}},
		{"test#29",
			args{attributeGlobalExtension12, []string{"test"}},
			AttributeMap{AttributeGlobalExtension12(): "test"}},
		{"test#30",
			args{attributeGlobalExtension13, []string{"test"}},
			AttributeMap{AttributeGlobalExtension13(): "test"}},
		{"test#31",
			args{attributeGlobalExtension14, []string{"test"}},
			AttributeMap{AttributeGlobalExtension14(): "test"}},
		{"test#32",
			args{attributeGlobalExtension15, []string{"test"}},
			AttributeMap{AttributeGlobalExtension15(): "test"}},
		{"test#33",
			args{attributeGlobalExtension16, []string{"test"}},
			AttributeMap{AttributeGlobalExtension16(): "test"}},
		{"test#34",
			args{attributeGlobalExtension17, []string{"test"}},
			AttributeMap{AttributeGlobalExtension17(): "test"}},
		{"test#35",
			args{attributeGlobalExtension18, []string{"test"}},
			AttributeMap{AttributeGlobalExtension18(): "test"}},
		{"test#36",
			args{attributeGlobalExtension19, []string{"test"}},
			AttributeMap{AttributeGlobalExtension19(): "test"}},
		{"test#37",
			args{attributeGlobalExtension20, []string{"test"}},
			AttributeMap{AttributeGlobalExtension20(): "test"}},
		{"test#38",
			args{attributeGlobalExtension21, []string{"test"}},
			AttributeMap{AttributeGlobalExtension21(): "test"}},
		{"test#39",
			args{attributeGlobalExtension22, []string{"test"}},
			AttributeMap{AttributeGlobalExtension22(): "test"}},
		{"test#40",
			args{attributeGlobalExtension23, []string{"test"}},
			AttributeMap{AttributeGlobalExtension23(): "test"}},
		{"test#41",
			args{attributeGlobalExtension24, []string{"test"}},
			AttributeMap{AttributeGlobalExtension24(): "test"}},
		{"test#42",
			args{attributeGlobalExtension25, []string{"test"}},
			AttributeMap{AttributeGlobalExtension25(): "test"}},
		{"test#43",
			args{attributeGlobalExtension26, []string{"test"}},
			AttributeMap{AttributeGlobalExtension26(): "test"}},
		{"test#44",
			args{attributeGlobalExtension27, []string{"test"}},
			AttributeMap{AttributeGlobalExtension27(): "test"}},
		{"test#45",
			args{attributeGlobalExtension28, []string{"test"}},
			AttributeMap{AttributeGlobalExtension28(): "test"}},
		{"test#46",
			args{attributeGlobalExtension29, []string{"test"}},
			AttributeMap{AttributeGlobalExtension29(): "test"}},
		{"test#47",
			args{attributeGlobalExtension30, []string{"test"}},
			AttributeMap{AttributeGlobalExtension30(): "test"}},
		{"test#48",
			args{attributeGroupType, []string{"10"}},
			AttributeMap{AttributeGroupType(): []string{"DISTRIBUTION", "GLOBAL", "UNIVERSAL"}}},
		{"test#49",
			args{attributeLocation, []string{"test"}},
			AttributeMap{AttributeLocation(): "test"}},
		{"test#50",
			args{attributeLastLogonTimestamp, []string{"128271382742968750"}},
			AttributeMap{AttributeLastLogonTimestamp(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#51",
			args{attributeMail, []string{"test"}},
			AttributeMap{AttributeMail(): "test"}},
		{"test#52",
			args{attributeMemberOf, []string{"test"}},
			AttributeMap{AttributeMemberOf(): []string{"test"}}},
		{"test#53",
			args{attributeMembers, []string{"test"}},
			AttributeMap{AttributeMembers(): []string{"test"}}},
		{"test#54",
			args{attributeMsRadiusFramedIpAddress, []string{"2130706433"}},
			AttributeMap{AttributeMsRadiusFramedIpAddress(): net.IP{127, 0, 0, 1}}},
		{"test#55",
			args{attributeName, []string{"test"}},
			AttributeMap{AttributeName(): "test"}},
		{"test#56",
			args{attributeObjectCategory, []string{"test"}},
			AttributeMap{AttributeObjectCategory(): "test"}},
		{"test#57",
			args{attributeObjectClass, []string{"test"}},
			AttributeMap{AttributeObjectClass(): []string{"test"}}},
		{"test#58",
			args{attributeObjectGUID, []string{"test"}},
			AttributeMap{AttributeObjectGUID(): `\x74\x65\x73\x74`}},
		{"test#59",
			args{attributeObjectSID, []string{"test"}},
			AttributeMap{AttributeObjectSID(): `\x74\x65\x73\x74`}},
		{"test#60",
			args{attributePasswordLastSet, []string{"128271382742968750"}},
			AttributeMap{AttributePasswordLastSet(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#61",
			args{attributePostalCode, []string{"12345"}},
			AttributeMap{AttributePostalCode(): "12345"}},
		{"test#60",
			args{attributeSamAccountName, []string{"test"}},
			AttributeMap{AttributeSamAccountName(): "test"}},
		{"test#61",
			args{attributeSamAccountType, []string{"805306368"}},
			AttributeMap{AttributeSamAccountType(): []string{"NORMAL_USER_ACCOUNT", "USER_OBJECT"}}},
		{"test#62",
			args{attributeSurname, []string{"test"}},
			AttributeMap{AttributeSurname(): "test"}},
		{"test#63",
			args{attributeStreetAddress, []string{"test"}},
			AttributeMap{AttributeStreetAddress(): "test"}},
		{"test#64",
			args{attributeUserAccountControl, []string{"8585216"}},
			AttributeMap{
				AttributeUserAccountControl():           []string{"DONT_EXPIRE_PASSWD", "MNS_LOGON_ACCOUNT", "PASSWORD_EXPIRED"},
				AttributeRaw("", "Enabled", TypeBool):   true,
				AttributeRaw("", "LockedOut", TypeBool): false,
			}},
		{"test#65",
			args{attributeUserCertificate, []string{"test"}},
			AttributeMap{AttributeUserCertificate(): `\x74\x65\x73\x74`}},
		{"test#66",
			args{attributeUserPrincipalName, []string{"test"}},
			AttributeMap{AttributeUserPrincipalName(): "test"}},
		{"test#67",
			args{attributeWhenChanged, []string{"128271382742968750"}},
			AttributeMap{AttributeWhenChanged(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#68",
			args{attributeWhenCreated, []string{"128271382742968750"}},
			AttributeMap{AttributeWhenCreated(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#69",
			args{AttributeRaw("unknown", "", TypeRaw), nil},
			AttributeMap{},
		},
		{"test#69",
			args{AttributeRaw("unknown", "", TypeRaw), []string{"test"}},
			AttributeMap{AttributeRaw("unknown", "", TypeRaw): "test"}},
		{"test#69",
			args{AttributeRaw("unknown", "", TypeRaw), []string{"test#1", "test#2"}},
			AttributeMap{AttributeRaw("unknown", "", TypeRaw): []string{"test#1", "test#2"}}},
		{"test#70",
			args{AttributeRaw("unknown", "", TypeBool), []string{"true"}},
			AttributeMap{AttributeRaw("unknown", "", TypeBool): true}},
		{"test#71",
			args{AttributeRaw("unknown", "", TypeBool), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeBool): []string{"invalid"}}},
		{"test#72",
			args{AttributeRaw("unknown", "", TypeDecimal), []string{"12.7"}},
			AttributeMap{AttributeRaw("unknown", "", TypeDecimal): float64(12.7)}},
		{"test#73",
			args{AttributeRaw("unknown", "", TypeDecimal), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeDecimal): []string{"invalid"}}},
		{"test#74",
			args{AttributeRaw("unknown", "", TypeGroupType), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeGroupType): []string{"invalid"}}},
		{"test#75",
			args{AttributeRaw("unknown", "", TypeInt), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeInt): []string{"invalid"}}},
		{"test#76",
			args{AttributeRaw("unknown", "", TypeIPv4Address), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeIPv4Address): []string{"invalid"}}},
		{"test#77",
			args{AttributeRaw("unknown", "", TypeSAMaccountType), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeSAMaccountType): []string{"invalid"}}},
		{"test#78",
			args{AttributeRaw("unknown", "", TypeTime), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeTime): []string{"invalid"}}},
		{"test#79",
			args{AttributeRaw("unknown", "", TypeUserAccountControl), []string{"invalid"}},
			AttributeMap{AttributeRaw("unknown", "", TypeUserAccountControl): []string{"invalid"}}},
		{"test#80",
			args{AttributeRaw("unknown", "", "unknown"), []string{"invalid"}},
			nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attr, attrMap, values := tt.args.a, AttributeMap{}, tt.args.v
			attr.Parse(values, &attrMap)

			changelogs, err := diff.Diff(attrMap, tt.want)
			if err != nil {
				t.Error(err)
			}

			for _, changelog := range changelogs {
				t.Errorf(`(%q).Parse(%v, ...) failed: Type: %q, Path: %v, From: %v, To: %v`, attr, values, changelog.Type, changelog.Path, changelog.From, changelog.To)
			}
		})
	}
}

func TestAttributeRegistryNotEmpty(t *testing.T) {
	if len(attributeRegistry) == 0 {
		t.Errorf(`no attributes registered`)
	}
}

func TestAttributesToAttributeList(t *testing.T) {
	for _, tt := range []struct {
		name string
		args Attributes
		want []string
	}{
		{"test#1", nil, nil},
		{"test#2",
			[]Attribute{AttributeAccountExpires(), AttributeUserAccountControl(), AttributeCommonName()},
			[]string{"accountexpires", "cn", "useraccountcontrol"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.ToAttributeList()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`AttributesToStringSlice(...) failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}

func TestLookupAttributeByLDAPDisplayName(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want *Attribute
	}{
		{"test#1", "name", &attributeName},
		{"test#2", "invalid", nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := LookupAttributeByLDAPDisplayName(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`LookupAttributeByLDAPDisplayName(%q) failed: got: %p, want: %p`, tt.args, got, tt.want)
			}
		})
	}
}
