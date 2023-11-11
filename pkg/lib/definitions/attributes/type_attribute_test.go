package attributes

import (
	"net"
	"reflect"
	"testing"
	"time"

	diff "github.com/r3labs/diff/v3"
)

func TestAppend(t *testing.T) {
	for _, tt := range []struct {
		name string
		args Attributes
		want Attributes
	}{
		{"test#1", Attributes{name, accountExpires, displayName, name}, Attributes{accountExpires, displayName, name}},
		{"test#2", Attributes{Any()}, Attributes{Any()}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got Attributes
			got.Append(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(Attributes).Append(...) failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		a Attribute
		v []string
	}
	for _, tt := range []struct {
		name string
		args args
		want Map
	}{
		{"test#1",
			args{accountExpires, []string{"128271382742968750"}},
			Map{AccountExpires(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#2",
			args{badPasswordTime, []string{"128271382742968750"}},
			Map{BadPasswordTime(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#3",
			args{badPasswordCount, []string{"3"}},
			Map{BadPasswordCount(): int64(3)}},
		{"test#4",
			args{commonName, []string{"test"}},
			Map{CommonName(): "test"}},
		{"test#5",
			args{countryCode, []string{"test"}},
			Map{CountryCode(): "test"}},
		{"test#6",
			args{countryName, []string{"test"}},
			Map{CountryName(): "test"}},
		{"test#7",
			args{countryNumber, []string{"103"}},
			Map{CountryNumber(): int64(103)}},
		{"test#8",
			args{company, []string{"test"}},
			Map{Company(): "test"}},
		{"test#9",
			args{department, []string{"test"}},
			Map{Department(): "test"}},
		{"test#10",
			args{departmentNumber, []string{"134"}},
			Map{DepartmentNumber(): "134"}},
		{"test#11",
			args{description, []string{"test"}},
			Map{Description(): "test"}},
		{"test#12",
			args{displayName, []string{"test"}},
			Map{DisplayName(): "test"}},
		{"test#13",
			args{distinguishedName, []string{"test"}},
			Map{DistinguishedName(): "test"}},
		{"test#14",
			args{division, []string{"test"}},
			Map{Division(): "test"}},
		{"test#15",
			args{dNSHostname, []string{"test"}},
			Map{DNSHostname(): "test"}},
		{"test#16",
			args{employeeID, []string{"test"}},
			Map{EmployeeID(): "test"}},
		{"test#17",
			args{givenName, []string{"test"}},
			Map{GivenName(): "test"}},
		{"test#18",
			args{globalExtension1, []string{"test"}},
			Map{GlobalExtension1(): "test"}},
		{"test#19",
			args{globalExtension2, []string{"test"}},
			Map{GlobalExtension2(): "test"}},
		{"test#20",
			args{globalExtension3, []string{"test"}},
			Map{GlobalExtension3(): "test"}},
		{"test#21",
			args{globalExtension4, []string{"test"}},
			Map{GlobalExtension4(): "test"}},
		{"test#22",
			args{globalExtension5, []string{"test"}},
			Map{GlobalExtension5(): "test"}},
		{"test#23",
			args{globalExtension6, []string{"test"}},
			Map{GlobalExtension6(): "test"}},
		{"test#24",
			args{globalExtension7, []string{"test"}},
			Map{GlobalExtension7(): "test"}},
		{"test#25",
			args{globalExtension8, []string{"test"}},
			Map{GlobalExtension8(): "test"}},
		{"test#26",
			args{globalExtension9, []string{"test"}},
			Map{GlobalExtension9(): "test"}},
		{"test#27",
			args{globalExtension10, []string{"test"}},
			Map{GlobalExtension10(): "test"}},
		{"test#28",
			args{globalExtension11, []string{"test"}},
			Map{GlobalExtension11(): "test"}},
		{"test#29",
			args{globalExtension12, []string{"test"}},
			Map{GlobalExtension12(): "test"}},
		{"test#30",
			args{globalExtension13, []string{"test"}},
			Map{GlobalExtension13(): "test"}},
		{"test#31",
			args{globalExtension14, []string{"test"}},
			Map{GlobalExtension14(): "test"}},
		{"test#32",
			args{globalExtension15, []string{"test"}},
			Map{GlobalExtension15(): "test"}},
		{"test#33",
			args{globalExtension16, []string{"test"}},
			Map{GlobalExtension16(): "test"}},
		{"test#34",
			args{globalExtension17, []string{"test"}},
			Map{GlobalExtension17(): "test"}},
		{"test#35",
			args{globalExtension18, []string{"test"}},
			Map{GlobalExtension18(): "test"}},
		{"test#36",
			args{globalExtension19, []string{"test"}},
			Map{GlobalExtension19(): "test"}},
		{"test#37",
			args{globalExtension20, []string{"test"}},
			Map{GlobalExtension20(): "test"}},
		{"test#38",
			args{globalExtension21, []string{"test"}},
			Map{GlobalExtension21(): "test"}},
		{"test#39",
			args{globalExtension22, []string{"test"}},
			Map{GlobalExtension22(): "test"}},
		{"test#40",
			args{globalExtension23, []string{"test"}},
			Map{GlobalExtension23(): "test"}},
		{"test#41",
			args{globalExtension24, []string{"test"}},
			Map{GlobalExtension24(): "test"}},
		{"test#42",
			args{globalExtension25, []string{"test"}},
			Map{GlobalExtension25(): "test"}},
		{"test#43",
			args{globalExtension26, []string{"test"}},
			Map{GlobalExtension26(): "test"}},
		{"test#44",
			args{globalExtension27, []string{"test"}},
			Map{GlobalExtension27(): "test"}},
		{"test#45",
			args{globalExtension28, []string{"test"}},
			Map{GlobalExtension28(): "test"}},
		{"test#46",
			args{globalExtension29, []string{"test"}},
			Map{GlobalExtension29(): "test"}},
		{"test#47",
			args{globalExtension30, []string{"test"}},
			Map{GlobalExtension30(): "test"}},
		{"test#48",
			args{groupType, []string{"10"}},
			Map{GroupType(): []string{"DISTRIBUTION", "GLOBAL", "UNIVERSAL"}}},
		{"test#49",
			args{location, []string{"test"}},
			Map{Location(): "test"}},
		{"test#50",
			args{lastLogonTimestamp, []string{"128271382742968750"}},
			Map{LastLogonTimestamp(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#51",
			args{mail, []string{"test"}},
			Map{Mail(): "test"}},
		{"test#52",
			args{memberOf, []string{"test"}},
			Map{MemberOf(): []string{"test"}}},
		{"test#53",
			args{members, []string{"test"}},
			Map{Members(): []string{"test"}}},
		{"test#54",
			args{msRadiusFramedIpAddress, []string{"2130706433"}},
			Map{MsRadiusFramedIpAddress(): net.IP{127, 0, 0, 1}}},
		{"test#55",
			args{name, []string{"test"}},
			Map{Name(): "test"}},
		{"test#56",
			args{objectCategory, []string{"test"}},
			Map{ObjectCategory(): "test"}},
		{"test#57",
			args{objectClass, []string{"test"}},
			Map{ObjectClass(): []string{"test"}}},
		{"test#58",
			args{objectGUID, []string{"test"}},
			Map{ObjectGUID(): `\x74\x65\x73\x74`}},
		{"test#59",
			args{objectSID, []string{"test"}},
			Map{ObjectSID(): `\x74\x65\x73\x74`}},
		{"test#60",
			args{passwordLastSet, []string{"128271382742968750"}},
			Map{PasswordLastSet(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#61",
			args{postalCode, []string{"12345"}},
			Map{PostalCode(): "12345"}},
		{"test#62",
			args{samAccountName, []string{"test"}},
			Map{SamAccountName(): "test"}},
		{"test#63",
			args{samAccountType, []string{"805306368"}},
			Map{SamAccountType(): []string{"NORMAL_USER_ACCOUNT", "USER_OBJECT"}}},
		{"test#64",
			args{surname, []string{"test"}},
			Map{Surname(): "test"}},
		{"test#65",
			args{streetAddress, []string{"test"}},
			Map{StreetAddress(): "test"}},
		{"test#66",
			args{userAccountControl, []string{"8585216"}},
			Map{
				UserAccountControl():           []string{"DONT_EXPIRE_PASSWD", "MNS_LOGON_ACCOUNT", "PASSWORD_EXPIRED"},
				Raw("", "Enabled", TypeBool):   true,
				Raw("", "LockedOut", TypeBool): false,
			}},
		{"test#67",
			args{userCertificate, []string{"test"}},
			Map{UserCertificate(): `\x74\x65\x73\x74`}},
		{"test#68",
			args{userPrincipalName, []string{"test"}},
			Map{UserPrincipalName(): "test"}},
		{"test#69",
			args{whenChanged, []string{"128271382742968750"}},
			Map{WhenChanged(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#70",
			args{whenCreated, []string{"128271382742968750"}},
			Map{WhenCreated(): time.Date(2007, 6, 24, 5, 57, 54, 296875000, time.UTC)}},
		{"test#71",
			args{Raw("unknown", "", TypeRaw), nil},
			Map{},
		},
		{"test#72",
			args{Raw("unknown", "", TypeRaw), []string{"test"}},
			Map{Raw("unknown", "", TypeRaw): "test"}},
		{"test#73",
			args{Raw("unknown", "", TypeRaw), []string{"test#1", "test#2"}},
			Map{Raw("unknown", "", TypeRaw): []string{"test#1", "test#2"}}},
		{"test#74",
			args{Raw("unknown", "", TypeBool), []string{"true"}},
			Map{Raw("unknown", "", TypeBool): true}},
		{"test#75",
			args{Raw("unknown", "", TypeBool), []string{"invalid"}},
			Map{Raw("unknown", "", TypeBool): []string{"invalid"}}},
		{"test#76",
			args{Raw("unknown", "", TypeDecimal), []string{"12.7"}},
			Map{Raw("unknown", "", TypeDecimal): float64(12.7)}},
		{"test#77",
			args{Raw("unknown", "", TypeDecimal), []string{"invalid"}},
			Map{Raw("unknown", "", TypeDecimal): []string{"invalid"}}},
		{"test#78",
			args{Raw("unknown", "", TypeGroupType), []string{"invalid"}},
			Map{Raw("unknown", "", TypeGroupType): []string{"invalid"}}},
		{"test#79",
			args{Raw("unknown", "", TypeInt), []string{"invalid"}},
			Map{Raw("unknown", "", TypeInt): []string{"invalid"}}},
		{"test#80",
			args{Raw("unknown", "", TypeIPv4Address), []string{"invalid"}},
			Map{Raw("unknown", "", TypeIPv4Address): []string{"invalid"}}},
		{"test#81",
			args{Raw("unknown", "", TypeSAMaccountType), []string{"invalid"}},
			Map{Raw("unknown", "", TypeSAMaccountType): []string{"invalid"}}},
		{"test#82",
			args{Raw("unknown", "", TypeTime), []string{"invalid"}},
			Map{Raw("unknown", "", TypeTime): []string{"invalid"}}},
		{"test#83",
			args{Raw("unknown", "", TypeUserAccountControl), []string{"invalid"}},
			Map{Raw("unknown", "", TypeUserAccountControl): []string{"invalid"}}},
		{"test#84",
			args{Raw("unknown", "", "unknown"), []string{"invalid"}},
			nil},
		{"test#85",
			args{UserPassword(), []string{"pass"}},
			Map{UserPassword(): "pass"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			attr, attrMap, values := tt.args.a, Map{}, tt.args.v
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

func TestToAttributeList(t *testing.T) {
	for _, tt := range []struct {
		name string
		args Attributes
		want []string
	}{
		{"test#1", nil, nil},
		{"test#2",
			[]Attribute{AccountExpires(), UserAccountControl(), CommonName(), CommonName(), Raw("", "custom", TypeRaw)},
			[]string{"AccountExpires", "CN", "Custom", "UserAccountControl"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.ToAttributeList()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(Attributes).ToStringSlice(...) failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}
