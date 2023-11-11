package attributes

import (
	"encoding/binary"
	"net"
	"slices"
	"strconv"
	"strings"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

var (
	attributeAccountExpires          = Attribute{"AccountExpires", "", TypeTime}.register()
	attributeBadPasswordTime         = Attribute{"BadPasswordTime", "", TypeTime}.register()
	attributeBadPasswordCount        = Attribute{"BadPwdCount", "BadPasswordCount", TypeInt}.register()
	attributeCommonName              = Attribute{"CN", "CommonName", TypeString}.register()
	attributeCountryCode             = Attribute{"C", "CountryCode", TypeString}.register()
	attributeCountryName             = Attribute{"CO", "CountryName", TypeString}.register()
	attributeCountryNumber           = Attribute{"CountryCode", "CountryNumber", TypeInt}.register()
	attributeCompany                 = Attribute{"Company", "", TypeString}.register()
	attributeDepartment              = Attribute{"Department", "", TypeString}.register()
	attributeDepartmentNumber        = Attribute{"DepartmentNumber", "", TypeString}.register()
	attributeDescription             = Attribute{"Description", "", TypeString}.register()
	attributeDisplayName             = Attribute{"DisplayName", "", TypeString}.register()
	attributeDistinguishedName       = Attribute{"DistinguishedName", "", TypeString}.register()
	attributeDivision                = Attribute{"Division", "", TypeString}.register()
	attributeDNSHostname             = Attribute{"DNSHostname", "", TypeString}.register()
	attributeEmployeeID              = Attribute{"EmployeeID", "", TypeString}.register()
	attributeGivenName               = Attribute{"GivenName", "", TypeString}.register()
	attributeGlobalExtension1        = Attribute{"Global-ExtensionAttribute1", "GlobalExtensionAttribute1", TypeString}.register()   // custom property
	attributeGlobalExtension2        = Attribute{"Global-ExtensionAttribute2", "GlobalExtensionAttribute2", TypeString}.register()   // custom property
	attributeGlobalExtension3        = Attribute{"Global-ExtensionAttribute3", "GlobalExtensionAttribute3", TypeString}.register()   // custom property
	attributeGlobalExtension4        = Attribute{"Global-ExtensionAttribute4", "GlobalExtensionAttribute4", TypeString}.register()   // custom property
	attributeGlobalExtension5        = Attribute{"Global-ExtensionAttribute5", "GlobalExtensionAttribute5", TypeString}.register()   // custom property
	attributeGlobalExtension6        = Attribute{"Global-ExtensionAttribute6", "GlobalExtensionAttribute6", TypeString}.register()   // custom property
	attributeGlobalExtension7        = Attribute{"Global-ExtensionAttribute7", "GlobalExtensionAttribute7", TypeString}.register()   // custom property
	attributeGlobalExtension8        = Attribute{"Global-ExtensionAttribute8", "GlobalExtensionAttribute8", TypeString}.register()   // custom property
	attributeGlobalExtension9        = Attribute{"Global-ExtensionAttribute9", "GlobalExtensionAttribute9", TypeString}.register()   // custom property
	attributeGlobalExtension10       = Attribute{"Global-ExtensionAttribute10", "GlobalExtensionAttribute10", TypeString}.register() // custom property
	attributeGlobalExtension11       = Attribute{"Global-ExtensionAttribute11", "GlobalExtensionAttribute11", TypeString}.register() // custom property
	attributeGlobalExtension12       = Attribute{"Global-ExtensionAttribute12", "GlobalExtensionAttribute12", TypeString}.register() // custom property
	attributeGlobalExtension13       = Attribute{"Global-ExtensionAttribute13", "GlobalExtensionAttribute13", TypeString}.register() // custom property
	attributeGlobalExtension14       = Attribute{"Global-ExtensionAttribute14", "GlobalExtensionAttribute14", TypeString}.register() // custom property
	attributeGlobalExtension15       = Attribute{"Global-ExtensionAttribute15", "GlobalExtensionAttribute15", TypeString}.register() // custom property
	attributeGlobalExtension16       = Attribute{"Global-ExtensionAttribute16", "GlobalExtensionAttribute16", TypeString}.register() // custom property
	attributeGlobalExtension17       = Attribute{"Global-ExtensionAttribute17", "GlobalExtensionAttribute17", TypeString}.register() // custom property
	attributeGlobalExtension18       = Attribute{"Global-ExtensionAttribute18", "GlobalExtensionAttribute18", TypeString}.register() // custom property
	attributeGlobalExtension19       = Attribute{"Global-ExtensionAttribute19", "GlobalExtensionAttribute19", TypeString}.register() // custom property
	attributeGlobalExtension20       = Attribute{"Global-ExtensionAttribute20", "GlobalExtensionAttribute20", TypeString}.register() // custom property
	attributeGlobalExtension21       = Attribute{"Global-ExtensionAttribute21", "GlobalExtensionAttribute21", TypeString}.register() // custom property
	attributeGlobalExtension22       = Attribute{"Global-ExtensionAttribute22", "GlobalExtensionAttribute22", TypeString}.register() // custom property
	attributeGlobalExtension23       = Attribute{"Global-ExtensionAttribute23", "GlobalExtensionAttribute23", TypeString}.register() // custom property
	attributeGlobalExtension24       = Attribute{"Global-ExtensionAttribute24", "GlobalExtensionAttribute24", TypeString}.register() // custom property
	attributeGlobalExtension25       = Attribute{"Global-ExtensionAttribute25", "GlobalExtensionAttribute25", TypeString}.register() // custom property
	attributeGlobalExtension26       = Attribute{"Global-ExtensionAttribute26", "GlobalExtensionAttribute26", TypeString}.register() // custom property
	attributeGlobalExtension27       = Attribute{"Global-ExtensionAttribute27", "GlobalExtensionAttribute27", TypeString}.register() // custom property
	attributeGlobalExtension28       = Attribute{"Global-ExtensionAttribute28", "GlobalExtensionAttribute28", TypeString}.register() // custom property
	attributeGlobalExtension29       = Attribute{"Global-ExtensionAttribute29", "GlobalExtensionAttribute29", TypeString}.register() // custom property
	attributeGlobalExtension30       = Attribute{"Global-ExtensionAttribute30", "GlobalExtensionAttribute30", TypeString}.register() // custom property
	attributeGroupType               = Attribute{"GroupType", "", TypeGroupType}.register()
	attributeLocation                = Attribute{"L", "Location", TypeString}.register()
	attributeLastLogonTimestamp      = Attribute{"LastLogonTimestamp", "", TypeTime}.register()
	attributeMail                    = Attribute{"Mail", "", TypeString}.register()
	attributeMemberOf                = Attribute{"MemberOf", "", TypeStringSlice}.register()
	attributeMembers                 = Attribute{"Member", "Members", TypeStringSlice}.register()
	attributeMsRadiusFramedIpAddress = Attribute{"MSRadiusFramedIPAddress", "", TypeIPv4Address} // custom property used in DMZ.register()
	attributeName                    = Attribute{"Name", "", TypeString}.register()
	attributeObjectCategory          = Attribute{"ObjectCategory", "", TypeString}.register()
	attributeObjectClass             = Attribute{"ObjectClass", "", TypeStringSlice}.register()
	attributeObjectGUID              = Attribute{"ObjectGuid", "", TypeHexString}.register()
	attributeObjectSID               = Attribute{"ObjectSid", "", TypeHexString}.register()
	attributePasswordLastSet         = Attribute{"PwdLastSet", "PasswordLastSet", TypeTime}.register()
	attributePostalCode              = Attribute{"PostalCode", "", TypeString}.register()
	attributeSamAccountName          = Attribute{"SAMAccountName", "", TypeString}.register()
	attributeSamAccountType          = Attribute{"SAMAccountType", "", TypeSAMaccountType}.register()
	attributeSurname                 = Attribute{"SN", "Surname", TypeString}.register()
	attributeStreetAddress           = Attribute{"StreetAddress", "", TypeString}.register()
	attributeUserAccountControl      = Attribute{"UserAccountControl", "", TypeUserAccountControl}.register()
	attributeUserCertificate         = Attribute{"UserCertificate", "", TypeHexString}.register()
	attributeUserPrincipalName       = Attribute{"UserPrincipalName", "", TypeString}.register()
	attributeWhenChanged             = Attribute{"WhenChanged", "", TypeTime}.register()
	attributeWhenCreated             = Attribute{"WhenCreated", "", TypeTime}.register()
)

var attributeRegistry Attributes

type Attribute struct {
	LDAPDisplayName string
	PrettyName      string
	Type            AttributeType
}

func (a Attribute) Parse(values []string, attributeMap *AttributeMap) {
	if len(values) == 0 || attributeMap == nil {
		return
	}

	switch a.Type {
	case TypeBool:
		parsed, err := strconv.ParseBool(values[0])
		if err == nil {
			(*attributeMap)[a] = parsed
		} else {
			(*attributeMap)[a] = values
		}

	case TypeDecimal:
		parsed, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			(*attributeMap)[a] = parsed
		} else {
			(*attributeMap)[a] = values
		}

	case TypeGroupType:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attributeMap)[a] = GroupType(parsed).Eval()
		} else {
			(*attributeMap)[a] = values
		}

	case TypeHexString:
		(*attributeMap)[a] = util.Hexify(values[0])

	case TypeInt:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attributeMap)[a] = parsed
		} else {
			(*attributeMap)[a] = values
		}

	case TypeIPv4Address:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, uint32(parsed))
			(*attributeMap)[a] = ip
		} else {
			(*attributeMap)[a] = values
		}

	case TypeRaw:
		if len(values) == 1 {
			(*attributeMap)[a] = values[0]
		} else {
			(*attributeMap)[a] = values
		}

	case TypeSAMaccountType:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attributeMap)[a] = SAMAccountType(parsed).Eval()
		} else {
			(*attributeMap)[a] = values
		}

	case TypeString:
		(*attributeMap)[a] = values[0]

	case TypeStringSlice:
		(*attributeMap)[a] = values

	case TypeTime:
		parsed, err := strconv.ParseInt(strings.Split(values[0], ".")[0], 10, 64)
		if err == nil {
			(*attributeMap)[a] = util.TimeAfter1601(parsed)
		} else {
			(*attributeMap)[a] = values
		}

	case TypeUserAccountControl:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			userAccountControl := UserAccountControl(parsed)
			(*attributeMap)[a] = userAccountControl.Eval()
			(*attributeMap)[AttributeRaw("", "Enabled", TypeBool).register()] = userAccountControl&USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE == 0
			(*attributeMap)[AttributeRaw("", "LockedOut", TypeBool).register()] = userAccountControl&USER_ACCOUNT_CONTROL_LOCKOUT != 0
		} else {
			(*attributeMap)[a] = values
		}

	default:
		return

	}
}

func (a Attribute) register() Attribute {
	attributeRegistry = append(attributeRegistry, a)
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
		list = append(list, strings.ToLower(attr.LDAPDisplayName))
	}

	slices.Sort(list)
	return list
}

func AttributeAccountExpires() Attribute          { return attributeAccountExpires }
func AttributeBadPasswordTime() Attribute         { return attributeBadPasswordTime }
func AttributeBadPasswordCount() Attribute        { return attributeBadPasswordCount }
func AttributeCommonName() Attribute              { return attributeCommonName }
func AttributeCountryCode() Attribute             { return attributeCountryCode }
func AttributeCountryName() Attribute             { return attributeCountryName }
func AttributeCountryNumber() Attribute           { return attributeCountryNumber }
func AttributeCompany() Attribute                 { return attributeCompany }
func AttributeDepartment() Attribute              { return attributeDepartment }
func AttributeDepartmentNumber() Attribute        { return attributeDepartmentNumber }
func AttributeDescription() Attribute             { return attributeDescription }
func AttributeDisplayName() Attribute             { return attributeDisplayName }
func AttributeDistinguishedName() Attribute       { return attributeDistinguishedName }
func AttributeDivision() Attribute                { return attributeDivision }
func AttributeDNSHostname() Attribute             { return attributeDNSHostname }
func AttributeEmployeeID() Attribute              { return attributeEmployeeID }
func AttributeGivenName() Attribute               { return attributeGivenName }
func AttributeGlobalExtension1() Attribute        { return attributeGlobalExtension1 }
func AttributeGlobalExtension2() Attribute        { return attributeGlobalExtension2 }
func AttributeGlobalExtension3() Attribute        { return attributeGlobalExtension3 }
func AttributeGlobalExtension4() Attribute        { return attributeGlobalExtension4 }
func AttributeGlobalExtension5() Attribute        { return attributeGlobalExtension5 }
func AttributeGlobalExtension6() Attribute        { return attributeGlobalExtension6 }
func AttributeGlobalExtension7() Attribute        { return attributeGlobalExtension7 }
func AttributeGlobalExtension8() Attribute        { return attributeGlobalExtension8 }
func AttributeGlobalExtension9() Attribute        { return attributeGlobalExtension9 }
func AttributeGlobalExtension10() Attribute       { return attributeGlobalExtension10 }
func AttributeGlobalExtension11() Attribute       { return attributeGlobalExtension11 }
func AttributeGlobalExtension12() Attribute       { return attributeGlobalExtension12 }
func AttributeGlobalExtension13() Attribute       { return attributeGlobalExtension13 }
func AttributeGlobalExtension14() Attribute       { return attributeGlobalExtension14 }
func AttributeGlobalExtension15() Attribute       { return attributeGlobalExtension15 }
func AttributeGlobalExtension16() Attribute       { return attributeGlobalExtension16 }
func AttributeGlobalExtension17() Attribute       { return attributeGlobalExtension17 }
func AttributeGlobalExtension18() Attribute       { return attributeGlobalExtension18 }
func AttributeGlobalExtension19() Attribute       { return attributeGlobalExtension19 }
func AttributeGlobalExtension20() Attribute       { return attributeGlobalExtension20 }
func AttributeGlobalExtension21() Attribute       { return attributeGlobalExtension21 }
func AttributeGlobalExtension22() Attribute       { return attributeGlobalExtension22 }
func AttributeGlobalExtension23() Attribute       { return attributeGlobalExtension23 }
func AttributeGlobalExtension24() Attribute       { return attributeGlobalExtension24 }
func AttributeGlobalExtension25() Attribute       { return attributeGlobalExtension25 }
func AttributeGlobalExtension26() Attribute       { return attributeGlobalExtension26 }
func AttributeGlobalExtension27() Attribute       { return attributeGlobalExtension27 }
func AttributeGlobalExtension28() Attribute       { return attributeGlobalExtension28 }
func AttributeGlobalExtension29() Attribute       { return attributeGlobalExtension29 }
func AttributeGlobalExtension30() Attribute       { return attributeGlobalExtension30 }
func AttributeGroupType() Attribute               { return attributeGroupType }
func AttributeLocation() Attribute                { return attributeLocation }
func AttributeLastLogonTimestamp() Attribute      { return attributeLastLogonTimestamp }
func AttributeMail() Attribute                    { return attributeMail }
func AttributeMemberOf() Attribute                { return attributeMemberOf }
func AttributeMembers() Attribute                 { return attributeMembers }
func AttributeMsRadiusFramedIpAddress() Attribute { return attributeMsRadiusFramedIpAddress }
func AttributeName() Attribute                    { return attributeName }
func AttributeObjectCategory() Attribute          { return attributeObjectCategory }
func AttributeObjectClass() Attribute             { return attributeObjectClass }
func AttributeObjectGUID() Attribute              { return attributeObjectGUID }
func AttributeObjectSID() Attribute               { return attributeObjectSID }
func AttributePasswordLastSet() Attribute         { return attributePasswordLastSet }
func AttributePostalCode() Attribute              { return attributePostalCode }
func AttributeSamAccountName() Attribute          { return attributeSamAccountName }
func AttributeSamAccountType() Attribute          { return attributeSamAccountType }
func AttributeSurname() Attribute                 { return attributeSurname }
func AttributeStreetAddress() Attribute           { return attributeStreetAddress }
func AttributeUserAccountControl() Attribute      { return attributeUserAccountControl }
func AttributeUserCertificate() Attribute         { return attributeUserCertificate }
func AttributeUserPrincipalName() Attribute       { return attributeUserPrincipalName }
func AttributeWhenChanged() Attribute             { return attributeWhenChanged }
func AttributeWhenCreated() Attribute             { return attributeWhenCreated }

func AttributeRaw(LDAPName, prettyName string, attrType AttributeType) Attribute {
	return Attribute{LDAPName, prettyName, attrType}
}

func LookupAttributeByLDAPDisplayName(in string) *Attribute {
	for _, attr := range attributeRegistry {
		switch {

		case // match either by LDAP display name or by pretty name
			strings.EqualFold(in, attr.LDAPDisplayName),
			attr.PrettyName != "" && strings.EqualFold(in, attr.PrettyName):

			return &attr

		}
	}

	return nil
}
