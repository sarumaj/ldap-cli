package attributes

import "strings"

var (
	accountExpires          = Attribute{"", "AccountExpires", "", TypeTime}.Register()
	badPasswordTime         = Attribute{"", "BadPasswordTime", "", TypeTime}.Register()
	badPasswordCount        = Attribute{"", "BadPwdCount", "BadPasswordCount", TypeInt}.Register()
	commonName              = Attribute{"", "CN", "CommonName", TypeString}.Register()
	countryCode             = Attribute{"", "C", "CountryCode", TypeString}.Register()
	countryName             = Attribute{"", "CO", "CountryName", TypeString}.Register()
	countryNumber           = Attribute{"", "CountryCode", "CountryNumber", TypeInt}.Register()
	company                 = Attribute{"", "Company", "", TypeString}.Register()
	department              = Attribute{"", "Department", "", TypeString}.Register()
	departmentNumber        = Attribute{"", "DepartmentNumber", "", TypeString}.Register()
	description             = Attribute{"", "Description", "", TypeString}.Register()
	displayName             = Attribute{"", "DisplayName", "", TypeString}.Register()
	distinguishedName       = Attribute{"DN", "DistinguishedName", "", TypeString}.Register()
	division                = Attribute{"", "Division", "", TypeString}.Register()
	dNSHostname             = Attribute{"", "DNSHostname", "", TypeString}.Register()
	employeeID              = Attribute{"", "EmployeeID", "", TypeString}.Register()
	givenName               = Attribute{"", "GivenName", "", TypeString}.Register()
	globalExtension1        = Attribute{"", "Global-ExtensionAttribute1", "GlobalExtensionAttribute1", TypeString}.Register()   // custom property
	globalExtension2        = Attribute{"", "Global-ExtensionAttribute2", "GlobalExtensionAttribute2", TypeString}.Register()   // custom property
	globalExtension3        = Attribute{"", "Global-ExtensionAttribute3", "GlobalExtensionAttribute3", TypeString}.Register()   // custom property
	globalExtension4        = Attribute{"", "Global-ExtensionAttribute4", "GlobalExtensionAttribute4", TypeString}.Register()   // custom property
	globalExtension5        = Attribute{"", "Global-ExtensionAttribute5", "GlobalExtensionAttribute5", TypeString}.Register()   // custom property
	globalExtension6        = Attribute{"", "Global-ExtensionAttribute6", "GlobalExtensionAttribute6", TypeString}.Register()   // custom property
	globalExtension7        = Attribute{"", "Global-ExtensionAttribute7", "GlobalExtensionAttribute7", TypeString}.Register()   // custom property
	globalExtension8        = Attribute{"", "Global-ExtensionAttribute8", "GlobalExtensionAttribute8", TypeString}.Register()   // custom property
	globalExtension9        = Attribute{"", "Global-ExtensionAttribute9", "GlobalExtensionAttribute9", TypeString}.Register()   // custom property
	globalExtension10       = Attribute{"", "Global-ExtensionAttribute10", "GlobalExtensionAttribute10", TypeString}.Register() // custom property
	globalExtension11       = Attribute{"", "Global-ExtensionAttribute11", "GlobalExtensionAttribute11", TypeString}.Register() // custom property
	globalExtension12       = Attribute{"", "Global-ExtensionAttribute12", "GlobalExtensionAttribute12", TypeString}.Register() // custom property
	globalExtension13       = Attribute{"", "Global-ExtensionAttribute13", "GlobalExtensionAttribute13", TypeString}.Register() // custom property
	globalExtension14       = Attribute{"", "Global-ExtensionAttribute14", "GlobalExtensionAttribute14", TypeString}.Register() // custom property
	globalExtension15       = Attribute{"", "Global-ExtensionAttribute15", "GlobalExtensionAttribute15", TypeString}.Register() // custom property
	globalExtension16       = Attribute{"", "Global-ExtensionAttribute16", "GlobalExtensionAttribute16", TypeString}.Register() // custom property
	globalExtension17       = Attribute{"", "Global-ExtensionAttribute17", "GlobalExtensionAttribute17", TypeString}.Register() // custom property
	globalExtension18       = Attribute{"", "Global-ExtensionAttribute18", "GlobalExtensionAttribute18", TypeString}.Register() // custom property
	globalExtension19       = Attribute{"", "Global-ExtensionAttribute19", "GlobalExtensionAttribute19", TypeString}.Register() // custom property
	globalExtension20       = Attribute{"", "Global-ExtensionAttribute20", "GlobalExtensionAttribute20", TypeString}.Register() // custom property
	globalExtension21       = Attribute{"", "Global-ExtensionAttribute21", "GlobalExtensionAttribute21", TypeString}.Register() // custom property
	globalExtension22       = Attribute{"", "Global-ExtensionAttribute22", "GlobalExtensionAttribute22", TypeString}.Register() // custom property
	globalExtension23       = Attribute{"", "Global-ExtensionAttribute23", "GlobalExtensionAttribute23", TypeString}.Register() // custom property
	globalExtension24       = Attribute{"", "Global-ExtensionAttribute24", "GlobalExtensionAttribute24", TypeString}.Register() // custom property
	globalExtension25       = Attribute{"", "Global-ExtensionAttribute25", "GlobalExtensionAttribute25", TypeString}.Register() // custom property
	globalExtension26       = Attribute{"", "Global-ExtensionAttribute26", "GlobalExtensionAttribute26", TypeString}.Register() // custom property
	globalExtension27       = Attribute{"", "Global-ExtensionAttribute27", "GlobalExtensionAttribute27", TypeString}.Register() // custom property
	globalExtension28       = Attribute{"", "Global-ExtensionAttribute28", "GlobalExtensionAttribute28", TypeString}.Register() // custom property
	globalExtension29       = Attribute{"", "Global-ExtensionAttribute29", "GlobalExtensionAttribute29", TypeString}.Register() // custom property
	globalExtension30       = Attribute{"", "Global-ExtensionAttribute30", "GlobalExtensionAttribute30", TypeString}.Register() // custom property
	groupType               = Attribute{"", "GroupType", "", TypeGroupType}.Register()
	location                = Attribute{"", "L", "Location", TypeString}.Register()
	lastLogonTimestamp      = Attribute{"", "LastLogonTimestamp", "", TypeTime}.Register()
	mail                    = Attribute{"", "Mail", "", TypeString}.Register()
	memberOf                = Attribute{"", "MemberOf", "", TypeStringSlice}.Register()
	members                 = Attribute{"", "Member", "Members", TypeStringSlice}.Register()
	msRadiusFramedIpAddress = Attribute{"", "MSRadiusFramedIPAddress", "", TypeIPv4Address} // custom property used in DMZ.Register()
	name                    = Attribute{"", "Name", "", TypeString}.Register()
	objectCategory          = Attribute{"", "ObjectCategory", "", TypeString}.Register()
	objectClass             = Attribute{"", "ObjectClass", "", TypeStringSlice}.Register()
	objectGUID              = Attribute{"", "ObjectGuid", "", TypeHexString}.Register()
	objectSID               = Attribute{"", "ObjectSid", "", TypeHexString}.Register()
	passwordLastSet         = Attribute{"", "PwdLastSet", "PasswordLastSet", TypeTime}.Register()
	postalCode              = Attribute{"", "PostalCode", "", TypeString}.Register()
	samAccountName          = Attribute{"", "SAMAccountName", "", TypeString}.Register()
	samAccountType          = Attribute{"", "SAMAccountType", "", TypeSAMaccountType}.Register()
	surname                 = Attribute{"", "SN", "Surname", TypeString}.Register()
	streetAddress           = Attribute{"", "StreetAddress", "", TypeString}.Register()
	unicodePassword         = Attribute{"", "UnicodePwd", "UnicodePassword", TypeString}.Register()
	userAccountControl      = Attribute{"", "UserAccountControl", "", TypeUserAccountControl}.Register()
	userCertificate         = Attribute{"", "UserCertificate", "", TypeHexString}.Register()
	userPrincipalName       = Attribute{"", "UserPrincipalName", "", TypeString}.Register()
	whenChanged             = Attribute{"", "WhenChanged", "", TypeTime}.Register()
	whenCreated             = Attribute{"", "WhenCreated", "", TypeTime}.Register()
)

// registry is a list of all known attributes
var registry Attributes

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
func UnicodePassword() Attribute         { return unicodePassword }
func UserAccountControl() Attribute      { return userAccountControl }
func UserCertificate() Attribute         { return userCertificate }
func UserPrincipalName() Attribute       { return userPrincipalName }
func WhenChanged() Attribute             { return whenChanged }
func WhenCreated() Attribute             { return whenCreated }

// Any returns an attribute that matches any attribute
func Any() Attribute { return Attribute{"", "*", "", TypeRaw} }

// Raw returns an attribute that matches the given LDAP name
func Raw(LDAPName, prettyName string, attrType Type) Attribute {
	return Attribute{"", LDAPName, prettyName, attrType}
}

// Lookup returns the attribute that matches the given LDAP name, pretty name or alias
func Lookup(in string) *Attribute {
	for _, attr := range registry {
		switch {

		case // match either by LDAP display name, pretty name or alias
			strings.EqualFold(in, attr.LDAPDisplayName),
			attr.PrettyName != "" && strings.EqualFold(in, attr.PrettyName),
			attr.Alias != "" && strings.EqualFold(in, attr.Alias):

			return &attr

		}
	}

	return nil
}

// LookupMany returns a list of attributes that match the given LDAP names, pretty names or aliases
// (for "*"" it returns all attributes)
func LookupMany(strict bool, in ...string) (list Attributes) {
	var asterisk bool
	for _, s := range in {
		if s == "*" {
			asterisk = true
			break
		}
	}

	if asterisk {
		if strict {
			list = append(list, registry...)
			list.Sort()
		} else {
			list = append(list, Any())
		}
		return
	}

	seen := make(map[Attribute]bool)
	for _, s := range in {
		if attr := Lookup(s); attr != nil && !seen[*attr] {
			list, seen[*attr] = append(list, *attr), true
		}
	}

	list.Sort()
	return list
}
