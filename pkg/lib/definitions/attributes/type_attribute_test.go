package attributes

import (
	"reflect"
	"testing"
	"time"

	"github.com/r3labs/diff/v3"
)

func TestAttributeParse(t *testing.T) {
	type args struct {
		a Attribute
		v []string
		m AttributeMap
	}
	for _, tt := range []struct {
		name string
		args args
		want AttributeMap
	}{
		{"test#1",
			args{attributeAccountExpires, []string{"9223372036854775807"}, make(AttributeMap)},
			AttributeMap{attributeAccountExpires: time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC)}},
		{"test#2",
			args{attributeBadPasswordTime, []string{"9223372036854775807"}, make(AttributeMap)},
			AttributeMap{attributeBadPasswordTime: time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC)}},
		{"test#3",
			args{attributeBadPasswordCount, []string{"3"}, make(AttributeMap)},
			AttributeMap{attributeBadPasswordCount: int64(3)}},
		{"test#4", args{attributeCommonName, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#5", args{attributeCountryCode, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#6", args{attributeCountryName, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#7", args{attributeCountryNumber, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#8", args{attributeCompany, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#9", args{attributeDepartment, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#10", args{attributeDepartmentNumber, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#11", args{attributeDescription, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#12", args{attributeDisplayName, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#13", args{attributeDistinguishedName, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#14", args{attributeDivision, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#15", args{attributeDNSHostname, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#16", args{attributeEmployeeID, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#17", args{attributeGivenName, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#18",
			args{attributeGlobalExtension1, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension1: "test"}},
		{"test#19",
			args{attributeGlobalExtension2, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension2: "test"}},
		{"test#20",
			args{attributeGlobalExtension3, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension3: "test"}},
		{"test#21",
			args{attributeGlobalExtension4, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension4: "test"}},
		{"test#22",
			args{attributeGlobalExtension5, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension5: "test"}},
		{"test#23",
			args{attributeGlobalExtension6, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension6: "test"}},
		{"test#24",
			args{attributeGlobalExtension7, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension7: "test"}},
		{"test#25",
			args{attributeGlobalExtension8, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension8: "test"}},
		{"test#26",
			args{attributeGlobalExtension9, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension9: "test"}},
		{"test#27",
			args{attributeGlobalExtension10, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension10: "test"}},
		{"test#28",
			args{attributeGlobalExtension11, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension11: "test"}},
		{"test#29",
			args{attributeGlobalExtension12, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension12: "test"}},
		{"test#30",
			args{attributeGlobalExtension13, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension13: "test"}},
		{"test#31",
			args{attributeGlobalExtension14, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension14: "test"}},
		{"test#32",
			args{attributeGlobalExtension15, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension15: "test"}},
		{"test#33",
			args{attributeGlobalExtension16, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension16: "test"}},
		{"test#34",
			args{attributeGlobalExtension17, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension17: "test"}},
		{"test#35",
			args{attributeGlobalExtension18, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension18: "test"}},
		{"test#36",
			args{attributeGlobalExtension19, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension19: "test"}},
		{"test#37",
			args{attributeGlobalExtension20, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension20: "test"}},
		{"test#38",
			args{attributeGlobalExtension21, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension21: "test"}},
		{"test#39",
			args{attributeGlobalExtension22, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension22: "test"}},
		{"test#40",
			args{attributeGlobalExtension23, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension23: "test"}},
		{"test#41",
			args{attributeGlobalExtension24, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension24: "test"}},
		{"test#42",
			args{attributeGlobalExtension25, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension25: "test"}},
		{"test#43",
			args{attributeGlobalExtension26, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension26: "test"}},
		{"test#44",
			args{attributeGlobalExtension27, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension27: "test"}},
		{"test#45",
			args{attributeGlobalExtension28, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension28: "test"}},
		{"test#46",
			args{attributeGlobalExtension29, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension29: "test"}},
		{"test#47",
			args{attributeGlobalExtension30, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeGlobalExtension30: "test"}},
		{"test#48", args{attributeGroupType, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#49", args{attributeLocation, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#50", args{attributeLastLogonTimestamp, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#51", args{attributeMail, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#52", args{attributeMemberOf, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#53", args{attributeMembers, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#54", args{attributeMsRadiusFramedIpAddress, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#55", args{attributeName, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#56", args{attributeObjectCategory, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#57", args{attributeObjectClass, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#58", args{attributeObjectGUID, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#59", args{attributeObjectSID, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#60", args{attributePasswordLastSet, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#61", args{attributePostalCode, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#60", args{attributeSamAccountName, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#61", args{attributeSamAccountType, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#62", args{attributeSurname, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#63", args{attributeStreetAddress, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#64", args{attributeUserAccountControl, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#65", args{attributeUserCertificate, nil, make(AttributeMap)}, AttributeMap{}},
		{"test#66",
			args{attributeUserPrincipalName, []string{"test"}, make(AttributeMap)},
			AttributeMap{attributeUserPrincipalName: "test"}},
		{"test#67",
			args{attributeWhenChanged, []string{"9223372036854775807"}, make(AttributeMap)},
			AttributeMap{attributeWhenChanged: time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC)}},
		{"test#68",
			args{attributeWhenCreated, []string{"9223372036854775807"}, make(AttributeMap)},
			AttributeMap{attributeWhenCreated: time.Date(2262, 4, 11, 23, 47, 16, 854775807, time.UTC)}},
		{"test#69",
			args{Attribute{"unknown", "", TypeRaw}, nil, make(AttributeMap)},
			AttributeMap{}},
		{"test#69",
			args{Attribute{"unknown", "", TypeRaw}, []string{"test"}, make(AttributeMap)},
			AttributeMap{Attribute{"unknown", "", TypeRaw}: "test"}},
		{"test#69",
			args{Attribute{"unknown", "", TypeRaw}, []string{"test#1", "test#2"}, make(AttributeMap)},
			AttributeMap{Attribute{"unknown", "", TypeRaw}: []string{"test#1", "test#2"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attr, attrMap, values := tt.args.a, tt.args.m, tt.args.v
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
