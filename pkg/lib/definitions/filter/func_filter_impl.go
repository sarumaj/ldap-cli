package filter

import (
	"fmt"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func HasNotExpired(strict bool) Filter {
	filter := Or(
		Filter{Attribute: attributes.AccountExpires(), Value: fmt.Sprint(0)},
		Filter{Attribute: attributes.AccountExpires(), Value: fmt.Sprint(1<<63 - 1)},
		Filter{
			Attribute: attributes.AccountExpires(),
			Value:     fmt.Sprintf(">=%d", libutil.TimeSince1601().Nanoseconds()/100),
		},
	)

	if strict {
		return And(filter, Filter{Attribute: attributes.AccountExpires(), Value: "*"})
	}

	return Or(filter, Not(Filter{Attribute: attributes.AccountExpires(), Value: "*"}))
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

func IsGroup() Filter { return Filter{attributes.ObjectClass(), "group", ""} }

func IsUser() Filter { return Filter{attributes.ObjectClass(), "user", ""} }
