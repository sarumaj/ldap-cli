package filter

import (
	"fmt"

	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func HasNotExpired(strict bool) Filter {
	filter := Or(
		Filter{Attribute: attributes.AttributeAccountExpires(), Value: fmt.Sprint(0)},
		Filter{Attribute: attributes.AttributeAccountExpires(), Value: fmt.Sprint(1<<63 - 1)},
		Filter{
			Attribute: attributes.AttributeAccountExpires(),
			Value:     fmt.Sprintf(">=%d", util.TimeSince1601().Nanoseconds()/100),
		},
	)

	if strict {
		return And(filter, Filter{Attribute: attributes.AttributeAccountExpires(), Value: "*"})
	}

	return Or(filter, Not(Filter{Attribute: attributes.AttributeAccountExpires(), Value: "*"}))
}

func IsDomainController() Filter {
	return And(
		Filter{attributes.AttributeObjectClass(), "computer", ""},
		Filter{attributes.AttributeUserAccountControl(), fmt.Sprintf("%d", attributes.USER_ACCOUNT_CONTROL_SERVER_TRUST_ACCOUNT), attributes.LDAP_MATCHING_RULE_BIT_AND},
	)
}

func IsEnabled() Filter {
	return Not(Filter{attributes.AttributeUserAccountControl(), "2", attributes.LDAP_MATCHING_RULE_BIT_AND})
}

func IsGroup() Filter { return Filter{attributes.AttributeObjectClass(), "group", ""} }

func IsUser() Filter { return Filter{attributes.AttributeObjectClass(), "user", ""} }
