package filter

import (
	"slices"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

type Alias struct{ Alias, Substitution string }

var aliases = []Alias{
	// filters
	{`${ENABLED}`, IsEnabled().String()},
	{`${DISABLED}`, Not(IsEnabled()).String()},
	{`${GROUP}`, IsGroup().String()},
	{`${USER}`, IsUser().String()},
	// matching rules
	{`${AND}`, string(attributes.LDAP_MATCHING_RULE_BIT_AND)},
	{`${OR}`, string(attributes.LDAP_MATCHING_RULE_BIT_OR)},
	{`${RECURSIVE}`, string(attributes.LDAP_MATCHING_RULE_IN_CHAIN)},
	{`${DATA}`, string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA)},
}

func ListAliases() []Alias {
	list := make([]Alias, len(aliases))
	_ = copy(list, aliases)

	slices.SortStableFunc(list, func(a, b Alias) int {
		if a.Alias > b.Alias {
			return 1
		}

		return -1
	})

	return list
}
