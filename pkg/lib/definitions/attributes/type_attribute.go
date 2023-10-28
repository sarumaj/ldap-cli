package attributes

import (
	"encoding/binary"
	"net"
	"slices"
	"strconv"
	"strings"

	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

var (
	accountExpires          = Attribute{"AccountExpires", "", TypeTime}.register()
	badPasswordTime         = Attribute{"BadPasswordTime", "", TypeTime}.register()
	badPasswordCount        = Attribute{"BadPwdCount", "BadPasswordCount", TypeInt}.register()
	commonName              = Attribute{"CN", "CommonName", TypeString}.register()
	countryCode             = Attribute{"C", "CountryCode", TypeString}.register()
	countryName             = Attribute{"CO", "CountryName", TypeString}.register()
	countryNumber           = Attribute{"CountryCode", "CountryNumber", TypeInt}.register()
	company                 = Attribute{"Company", "", TypeString}.register()
	department              = Attribute{"Department", "", TypeString}.register()
	departmentNumber        = Attribute{"DepartmentNumber", "", TypeString}.register()
	description             = Attribute{"Description", "", TypeString}.register()
	displayName             = Attribute{"DisplayName", "", TypeString}.register()
	distinguishedName       = Attribute{"DistinguishedName", "", TypeString}.register()
	division                = Attribute{"Division", "", TypeString}.register()
	dNSHostname             = Attribute{"DNSHostname", "", TypeString}.register()
	employeeID              = Attribute{"EmployeeID", "", TypeString}.register()
	givenName               = Attribute{"GivenName", "", TypeString}.register()
	globalExtension1        = Attribute{"Global-ExtensionAttribute1", "GlobalExtensionAttribute1", TypeString}.register()   // custom property
	globalExtension2        = Attribute{"Global-ExtensionAttribute2", "GlobalExtensionAttribute2", TypeString}.register()   // custom property
	globalExtension3        = Attribute{"Global-ExtensionAttribute3", "GlobalExtensionAttribute3", TypeString}.register()   // custom property
	globalExtension4        = Attribute{"Global-ExtensionAttribute4", "GlobalExtensionAttribute4", TypeString}.register()   // custom property
	globalExtension5        = Attribute{"Global-ExtensionAttribute5", "GlobalExtensionAttribute5", TypeString}.register()   // custom property
	globalExtension6        = Attribute{"Global-ExtensionAttribute6", "GlobalExtensionAttribute6", TypeString}.register()   // custom property
	globalExtension7        = Attribute{"Global-ExtensionAttribute7", "GlobalExtensionAttribute7", TypeString}.register()   // custom property
	globalExtension8        = Attribute{"Global-ExtensionAttribute8", "GlobalExtensionAttribute8", TypeString}.register()   // custom property
	globalExtension9        = Attribute{"Global-ExtensionAttribute9", "GlobalExtensionAttribute9", TypeString}.register()   // custom property
	globalExtension10       = Attribute{"Global-ExtensionAttribute10", "GlobalExtensionAttribute10", TypeString}.register() // custom property
	globalExtension11       = Attribute{"Global-ExtensionAttribute11", "GlobalExtensionAttribute11", TypeString}.register() // custom property
	globalExtension12       = Attribute{"Global-ExtensionAttribute12", "GlobalExtensionAttribute12", TypeString}.register() // custom property
	globalExtension13       = Attribute{"Global-ExtensionAttribute13", "GlobalExtensionAttribute13", TypeString}.register() // custom property
	globalExtension14       = Attribute{"Global-ExtensionAttribute14", "GlobalExtensionAttribute14", TypeString}.register() // custom property
	globalExtension15       = Attribute{"Global-ExtensionAttribute15", "GlobalExtensionAttribute15", TypeString}.register() // custom property
	globalExtension16       = Attribute{"Global-ExtensionAttribute16", "GlobalExtensionAttribute16", TypeString}.register() // custom property
	globalExtension17       = Attribute{"Global-ExtensionAttribute17", "GlobalExtensionAttribute17", TypeString}.register() // custom property
	globalExtension18       = Attribute{"Global-ExtensionAttribute18", "GlobalExtensionAttribute18", TypeString}.register() // custom property
	globalExtension19       = Attribute{"Global-ExtensionAttribute19", "GlobalExtensionAttribute19", TypeString}.register() // custom property
	globalExtension20       = Attribute{"Global-ExtensionAttribute20", "GlobalExtensionAttribute20", TypeString}.register() // custom property
	globalExtension21       = Attribute{"Global-ExtensionAttribute21", "GlobalExtensionAttribute21", TypeString}.register() // custom property
	globalExtension22       = Attribute{"Global-ExtensionAttribute22", "GlobalExtensionAttribute22", TypeString}.register() // custom property
	globalExtension23       = Attribute{"Global-ExtensionAttribute23", "GlobalExtensionAttribute23", TypeString}.register() // custom property
	globalExtension24       = Attribute{"Global-ExtensionAttribute24", "GlobalExtensionAttribute24", TypeString}.register() // custom property
	globalExtension25       = Attribute{"Global-ExtensionAttribute25", "GlobalExtensionAttribute25", TypeString}.register() // custom property
	globalExtension26       = Attribute{"Global-ExtensionAttribute26", "GlobalExtensionAttribute26", TypeString}.register() // custom property
	globalExtension27       = Attribute{"Global-ExtensionAttribute27", "GlobalExtensionAttribute27", TypeString}.register() // custom property
	globalExtension28       = Attribute{"Global-ExtensionAttribute28", "GlobalExtensionAttribute28", TypeString}.register() // custom property
	globalExtension29       = Attribute{"Global-ExtensionAttribute29", "GlobalExtensionAttribute29", TypeString}.register() // custom property
	globalExtension30       = Attribute{"Global-ExtensionAttribute30", "GlobalExtensionAttribute30", TypeString}.register() // custom property
	groupType               = Attribute{"GroupType", "", TypeGroupType}.register()
	location                = Attribute{"L", "Location", TypeString}.register()
	lastLogonTimestamp      = Attribute{"LastLogonTimestamp", "", TypeTime}.register()
	mail                    = Attribute{"Mail", "", TypeString}.register()
	memberOf                = Attribute{"MemberOf", "", TypeStringSlice}.register()
	members                 = Attribute{"Member", "Members", TypeStringSlice}.register()
	msRadiusFramedIpAddress = Attribute{"MSRadiusFramedIPAddress", "", TypeIPv4Address} // custom property used in DMZ.register()
	name                    = Attribute{"Name", "", TypeString}.register()
	objectCategory          = Attribute{"ObjectCategory", "", TypeString}.register()
	objectClass             = Attribute{"ObjectClass", "", TypeStringSlice}.register()
	objectGUID              = Attribute{"ObjectGuid", "", TypeHexString}.register()
	objectSID               = Attribute{"ObjectSid", "", TypeHexString}.register()
	passwordLastSet         = Attribute{"PwdLastSet", "PasswordLastSet", TypeTime}.register()
	postalCode              = Attribute{"PostalCode", "", TypeString}.register()
	samAccountName          = Attribute{"SAMAccountName", "", TypeString}.register()
	samAccountType          = Attribute{"SAMAccountType", "", TypeSAMaccountType}.register()
	surname                 = Attribute{"SN", "Surname", TypeString}.register()
	streetAddress           = Attribute{"StreetAddress", "", TypeString}.register()
	userAccountControl      = Attribute{"UserAccountControl", "", TypeUserAccountControl}.register()
	userCertificate         = Attribute{"UserCertificate", "", TypeHexString}.register()
	userPrincipalName       = Attribute{"UserPrincipalName", "", TypeString}.register()
	whenChanged             = Attribute{"WhenChanged", "", TypeTime}.register()
	whenCreated             = Attribute{"WhenCreated", "", TypeTime}.register()
)

var registry Attributes

type Attribute struct {
	LDAPDisplayName string
	PrettyName      string
	Type            Type
}

func (a Attribute) Parse(values []string, attrMap *Map) {
	if len(values) == 0 || attrMap == nil {
		return
	}

	switch a.Type {
	case TypeBool:
		parsed, err := strconv.ParseBool(values[0])
		if err == nil {
			(*attrMap)[a] = parsed
		} else {
			(*attrMap)[a] = values
		}

	case TypeDecimal:
		parsed, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			(*attrMap)[a] = parsed
		} else {
			(*attrMap)[a] = values
		}

	case TypeGroupType:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = FlagsetGroupType(parsed).Eval()
		} else {
			(*attrMap)[a] = values
		}

	case TypeHexString:
		(*attrMap)[a] = libutil.Hexify(values[0])

	case TypeInt:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = parsed
		} else {
			(*attrMap)[a] = values
		}

	case TypeIPv4Address:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, uint32(parsed))
			(*attrMap)[a] = ip
		} else {
			(*attrMap)[a] = values
		}

	case TypeRaw:
		if len(values) == 1 {
			(*attrMap)[a] = values[0]
		} else {
			(*attrMap)[a] = values
		}

	case TypeSAMaccountType:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = FlagSAMAccountType(parsed).Eval()
		} else {
			(*attrMap)[a] = values
		}

	case TypeString:
		(*attrMap)[a] = values[0]

	case TypeStringSlice:
		(*attrMap)[a] = values

	case TypeTime:
		parsed, err := strconv.ParseInt(strings.Split(values[0], ".")[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = libutil.TimeAfter1601(parsed)
		} else {
			(*attrMap)[a] = values
		}

	case TypeUserAccountControl:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			userAccountControl := FlagsetUserAccountControl(parsed)
			(*attrMap)[a] = userAccountControl.Eval()
			(*attrMap)[Raw("", "Enabled", TypeBool).register()] = userAccountControl&USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE == 0
			(*attrMap)[Raw("", "LockedOut", TypeBool).register()] = userAccountControl&USER_ACCOUNT_CONTROL_LOCKOUT != 0
		} else {
			(*attrMap)[a] = values
		}

	default:
		return

	}
}

func (a Attribute) register() Attribute {
	registry = append(registry, a)
	return a
}

func (a Attribute) String() string {
	if a.PrettyName != "" {
		return a.PrettyName
	}

	return a.LDAPDisplayName
}

type Attributes []Attribute

func (a Attributes) ToAttributeList() []string {
	var list []string
	for _, attr := range a {
		if attr.LDAPDisplayName != "" {
			list = append(list, strings.ToLower(attr.LDAPDisplayName))
		} else {
			list = append(list, strings.ToLower(attr.PrettyName))
		}
	}

	slices.Sort(list)
	return list
}

func AccountExpires() Attribute          { return accountExpires }
func BadPasswordTime() Attribute         { return badPasswordTime }
func BadPasswordCount() Attribute        { return badPasswordCount }
func CommonName() Attribute              { return commonName }
func CountryCode() Attribute             { return countryCode }
func CountryName() Attribute             { return countryName }
func CountryNumber() Attribute           { return countryNumber }
func Company() Attribute                 { return company }
func Department() Attribute              { return department }
func DepartmentNumber() Attribute        { return departmentNumber }
func Description() Attribute             { return description }
func DisplayName() Attribute             { return displayName }
func DistinguishedName() Attribute       { return distinguishedName }
func Division() Attribute                { return division }
func DNSHostname() Attribute             { return dNSHostname }
func EmployeeID() Attribute              { return employeeID }
func GivenName() Attribute               { return givenName }
func GlobalExtension1() Attribute        { return globalExtension1 }
func GlobalExtension2() Attribute        { return globalExtension2 }
func GlobalExtension3() Attribute        { return globalExtension3 }
func GlobalExtension4() Attribute        { return globalExtension4 }
func GlobalExtension5() Attribute        { return globalExtension5 }
func GlobalExtension6() Attribute        { return globalExtension6 }
func GlobalExtension7() Attribute        { return globalExtension7 }
func GlobalExtension8() Attribute        { return globalExtension8 }
func GlobalExtension9() Attribute        { return globalExtension9 }
func GlobalExtension10() Attribute       { return globalExtension10 }
func GlobalExtension11() Attribute       { return globalExtension11 }
func GlobalExtension12() Attribute       { return globalExtension12 }
func GlobalExtension13() Attribute       { return globalExtension13 }
func GlobalExtension14() Attribute       { return globalExtension14 }
func GlobalExtension15() Attribute       { return globalExtension15 }
func GlobalExtension16() Attribute       { return globalExtension16 }
func GlobalExtension17() Attribute       { return globalExtension17 }
func GlobalExtension18() Attribute       { return globalExtension18 }
func GlobalExtension19() Attribute       { return globalExtension19 }
func GlobalExtension20() Attribute       { return globalExtension20 }
func GlobalExtension21() Attribute       { return globalExtension21 }
func GlobalExtension22() Attribute       { return globalExtension22 }
func GlobalExtension23() Attribute       { return globalExtension23 }
func GlobalExtension24() Attribute       { return globalExtension24 }
func GlobalExtension25() Attribute       { return globalExtension25 }
func GlobalExtension26() Attribute       { return globalExtension26 }
func GlobalExtension27() Attribute       { return globalExtension27 }
func GlobalExtension28() Attribute       { return globalExtension28 }
func GlobalExtension29() Attribute       { return globalExtension29 }
func GlobalExtension30() Attribute       { return globalExtension30 }
func GroupType() Attribute               { return groupType }
func Location() Attribute                { return location }
func LastLogonTimestamp() Attribute      { return lastLogonTimestamp }
func Mail() Attribute                    { return mail }
func MemberOf() Attribute                { return memberOf }
func Members() Attribute                 { return members }
func MsRadiusFramedIpAddress() Attribute { return msRadiusFramedIpAddress }
func Name() Attribute                    { return name }
func ObjectCategory() Attribute          { return objectCategory }
func ObjectClass() Attribute             { return objectClass }
func ObjectGUID() Attribute              { return objectGUID }
func ObjectSID() Attribute               { return objectSID }
func PasswordLastSet() Attribute         { return passwordLastSet }
func PostalCode() Attribute              { return postalCode }
func SamAccountName() Attribute          { return samAccountName }
func SamAccountType() Attribute          { return samAccountType }
func Surname() Attribute                 { return surname }
func StreetAddress() Attribute           { return streetAddress }
func UserAccountControl() Attribute      { return userAccountControl }
func UserCertificate() Attribute         { return userCertificate }
func UserPrincipalName() Attribute       { return userPrincipalName }
func WhenChanged() Attribute             { return whenChanged }
func WhenCreated() Attribute             { return whenCreated }

func Raw(LDAPName, prettyName string, attrType Type) Attribute {
	return Attribute{LDAPName, prettyName, attrType}
}

func Lookup(in string) *Attribute {
	for _, attr := range registry {
		switch {

		case // match either by LDAP display name or by pretty name
			strings.EqualFold(in, attr.LDAPDisplayName),
			attr.PrettyName != "" && strings.EqualFold(in, attr.PrettyName):

			return &attr

		}
	}

	return nil
}
