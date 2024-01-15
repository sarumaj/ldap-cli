package filter

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

var aliases = []Alias{
	// filters
	{`$ENABLED`, nil, func([]string) string { return IsEnabled().String() }},
	{`$DISABLED`, nil, func([]string) string { return Not(IsEnabled()).String() }},
	{`$GROUP`, nil, func([]string) string { return IsGroup().String() }},
	{`$USER`, nil, func([]string) string { return IsUser().String() }},
	{`$DC`, nil, func([]string) string { return IsDomainController().String() }},
	{`$EXPIRED`, nil, func([]string) string { return HasExpired().String() }},
	{`$NOT_EXPIRED`, nil, func([]string) string { return Not(HasExpired()).String() }},
	{`$ID`, []string{"id"}, func(params []string) string { return ByID(params[0]).String() }},
	{`$MEMBER_OF`, []string{"dn"}, func(params []string) string { return MemberOf(params[0], false).String() }},
	{`$MEMBER_OF_RECURSIVE`, []string{"dn"}, func(params []string) string { return MemberOf(params[0], true).String() }},
	// matching rules
	{`$AND`, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_AND) }},
	{`$OR`, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_OR) }},
	{`$RECURSIVE`, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_IN_CHAIN) }},
	{`$DATA`, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA) }},
}

type Alias struct {
	ID           string
	Parameters   []string
	Substitution func([]string) string
}

func (a Alias) String() string {
	if len(a.Parameters) > 0 {
		return fmt.Sprintf("%s(%s)", a.ID, strings.Join(a.Parameters, "; "))
	}

	return a.ID
}

func ListAliases() []Alias {
	list := make([]Alias, len(aliases))
	_ = copy(list, aliases)

	slices.SortStableFunc(list, func(a, b Alias) int {
		if a.ID > b.ID {
			return 1
		}

		return -1
	})

	return list
}

func ReplaceAliases(raw string) string {
	for _, alias := range aliases {
		if len(alias.Parameters) > 0 {
			for _, match := range regexp.
				MustCompile(`(?P<Alias>`+regexp.QuoteMeta(alias.ID)+`)\((?P<Parameters>[^\)]+)\)`).
				FindAllStringSubmatch(raw, -1) {

				raw = strings.ReplaceAll(raw, match[1]+"("+match[2]+")", alias.Substitution(regexp.MustCompile(`;\s*`).Split(match[2], -1)))
			}

		} else {
			raw = strings.ReplaceAll(raw, alias.ID, alias.Substitution(nil))
		}
	}

	return raw
}
