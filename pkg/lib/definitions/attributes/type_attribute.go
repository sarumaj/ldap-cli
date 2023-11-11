package attributes

import "slices"

const (
	AttributeAccountExpires          Attribute = "accountExpires"
	AttributeBadPasswordTime         Attribute = "badPasswordTime"
	AttributeBadPasswordCount        Attribute = "badPwdCount"
	AttributeCommonName              Attribute = "cn"
	AttributeCountryCode             Attribute = "c"
	AttributeCountryName             Attribute = "co"
	AttributeDescription             Attribute = "description"
	AttributeDisplayName             Attribute = "displayName"
	AttributeDistinguishedName       Attribute = "distinguishedName"
	AttributeEnabled                 Attribute = "userAccountControl"
	AttributeGivenName               Attribute = "givenName"
	AttributeGlobalExtension11       Attribute = "global-ExtensionAttribute11"
	AttributeGlobalExtension22       Attribute = "global-ExtensionAttribute22"
	AttributeGlobalExtension26       Attribute = "global-ExtensionAttribute26"
	AttributeGroupType               Attribute = "groupType"
	AttributeHostname                Attribute = "dnsHostname"
	AttributeMail                    Attribute = "mail"
	AttributeMemberOf                Attribute = "memberOf"
	AttributeMembers                 Attribute = "member"
	AttributeName                    Attribute = "name"
	AttributeOrganizationalUnit      Attribute = "ou"
	AttributeObjectCategory          Attribute = "objectCategory"
	AttributeObjectClass             Attribute = "objectClass"
	AttributeObjectGUID              Attribute = "objectGuid"
	AttributeSamAccountName          Attribute = "sAMAccountName"
	AttributeSamAccountType          Attribute = "sAMAccountType"
	AttributeSID                     Attribute = "objectSid"
	AttributeSurname                 Attribute = "sn"
	AttributeUserAccountControl      Attribute = "userAccountControl"
	AttributeUserPrincipalName       Attribute = "userPrincipalName"
	AttributeMsRadiusFramedIpAddress Attribute = "msRadiusFramedIPAddress" // custom property used in DMZ
)

type Attribute string

func (s Attribute) String() string { return string(s) }

func AttributesToStringSlice(attrs ...Attribute) (list []string) {
	for _, attr := range attrs {
		list = append(list, attr.String())
	}

	slices.Sort(list)
	return list
}
