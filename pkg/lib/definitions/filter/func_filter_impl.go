package filter

import (
	"fmt"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func ByID(id string) Filter {
	return Or(
		Filter{Attribute: attributes.CommonName(), Value: id},
		Filter{Attribute: attributes.DisplayName(), Value: id},
		Filter{Attribute: attributes.DistinguishedName(), Value: id}.ExpandAlias(),
		Filter{Attribute: attributes.Name(), Value: id},
		Filter{Attribute: attributes.SamAccountName(), Value: id},
		Filter{Attribute: attributes.UserPrincipalName(), Value: id},
	)
}

func HasExpired() Filter {
	return And(
		Filter{Attribute: attributes.AccountExpires(), Value: ">0"},
		Filter{Attribute: attributes.AccountExpires(), Value: fmt.Sprintf("<%d", int64(1<<63-1))},
		Filter{
			Attribute: attributes.AccountExpires(),
			Value:     fmt.Sprintf("<%d", libutil.TimeSince1601().Nanoseconds()/100),
		},
		Filter{Attribute: attributes.AccountExpires(), Value: "*"},
	)
}

func IsDomainController() Filter {
	return And(
		Filter{attributes.ObjectClass(), "computer", ""},
		Filter{attributes.UserAccountControl(), fmt.Sprintf("%d", attributes.USER_ACCOUNT_CONTROL_SERVER_TRUST_ACCOUNT), attributes.LDAP_MATCHING_RULE_BIT_AND},
	)
}

func IsEnabled() Filter {
	return Not(Filter{attributes.UserAccountControl(), "2", attributes.LDAP_MATCHING_RULE_BIT_AND})
}

func IsGroup() Filter {
	return Or(
		Filter{attributes.ObjectClass(), "group", ""},
		Filter{attributes.ObjectClass(), "posixGroup", ""}, // for testing purposes against openldap
	)
}

func IsUser() Filter {
	return Or(
		Filter{attributes.ObjectClass(), "user", ""},
		Filter{attributes.ObjectClass(), "posixAccount", ""}, // for testing purposes against openldap
	)
}

func MemberOf(parent string, recursive bool) Filter {
	if recursive {
		return Filter{Attribute: attributes.MemberOf(), Value: parent, Rule: attributes.LDAP_MATCHING_RULE_IN_CHAIN}
	}

	return Filter{Attribute: attributes.MemberOf(), Value: parent}
}
