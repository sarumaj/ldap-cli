package filter

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

// list of aliases
var aliases = []Alias{
	// filters
	{`$ENABLED`, FilterImplementation, nil, func([]string) string { return IsEnabled().String() }},
	{`$DISABLED`, FilterImplementation, nil, func([]string) string { return Not(IsEnabled()).String() }},
	{`$GROUP`, FilterImplementation, nil, func([]string) string { return IsGroup().String() }},
	{`$USER`, FilterImplementation, nil, func([]string) string { return IsUser().String() }},
	{`$DC`, FilterImplementation, nil, func([]string) string { return IsDomainController().String() }},
	{`$EXPIRED`, FilterImplementation, nil, func([]string) string { return HasExpired().String() }},
	{`$NOT_EXPIRED`, FilterImplementation, nil, func([]string) string { return Not(HasExpired()).String() }},
	{`$ID`, FilterImplementation, []string{"id"}, func(params []string) string { return ByID(params[0]).String() }},
	{`$MEMBER_OF`, FilterImplementation, []string{"dn"}, func(params []string) string { return MemberOf(params[0], false).String() }},
	{`$NESTED_MEMBER_OF`, FilterImplementation, []string{"dn"}, func(params []string) string { return MemberOf(params[0], true).String() }},
	// matching rules
	{`$BAND`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_AND) }},
	{`$BOR`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_OR) }},
	{`$RECURSIVE`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_IN_CHAIN) }},
	{`$DATA`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA) }},
}

const (
	// FilterImplementation is the kind of an alias that is a custom filter implementation
	FilterImplementation Kind = "filter implementation"
	// MatchingRule is the kind of an alias that is a matching rule bit mask
	MatchingRule Kind = "matching rule"
)

// delimiter is used to split alias parameters
var delimiter = regexp.MustCompile(`;\s*`)

// Alias is used to define shortcuts for filters and matching rules
type Alias struct {
	// ID is the identifier of the alias
	ID string
	// Kind is the kind of the alias
	Kind Kind
	// Parameters are the parameters of the alias used for substitution
	Parameters []string
	// Substitution is the function used to substitute the alias
	Substitution func([]string) string
}

// String returns a string representation of an alias
func (a Alias) String() string {
	if len(a.Parameters) > 0 {
		return fmt.Sprintf("%s(%s)", a.ID, strings.Join(a.Parameters, "; "))
	}

	return a.ID
}

// Kind is used to define the kind of an alias
type Kind string

// ListAliases returns a list of all aliases
func ListAliases() []Alias {
	list := make([]Alias, len(aliases))
	_ = copy(list, aliases)

	slices.SortStableFunc(list, func(a, b Alias) int {
		if string(a.Kind) == string(b.Kind) && a.ID > b.ID {
			return 1
		}

		return -1
	})

	return list
}

// ReplaceAliases replaces aliases in a raw string
func ReplaceAliases(raw string) string {
	for _, alias := range aliases {
		if len(alias.Parameters) > 0 {
			for _, match := range regexp.
				MustCompile(`(?P<Alias>`+regexp.QuoteMeta(alias.ID)+`)\((?P<Parameters>[^\)]+)\)`).
				FindAllStringSubmatch(raw, -1) {

				raw = strings.ReplaceAll(raw, match[1]+"("+match[2]+")", alias.Substitution(delimiter.Split(match[2], -1)))
			}

		} else {
			raw = strings.ReplaceAll(raw, alias.ID, alias.Substitution(nil))
		}
	}

	return raw
}
