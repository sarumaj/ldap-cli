package filter

import (
	"strconv"
	"strings"

	attributes "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/attributes"
)

// list of registry
var registry []Alias

// filters
var enabled = Alias{`$ENABLED`, FilterImplementation, nil, func([]string) string { return IsEnabled().String() }}.Register()
var disabled = Alias{`$DISABLED`, FilterImplementation, nil, func([]string) string { return Not(IsEnabled()).String() }}.Register()
var group = Alias{`$GROUP`, FilterImplementation, nil, func([]string) string { return IsGroup().String() }}.Register()
var user = Alias{`$USER`, FilterImplementation, nil, func([]string) string { return IsUser().String() }}.Register()
var dc = Alias{`$DC`, FilterImplementation, nil, func([]string) string { return IsDomainController().String() }}.Register()
var expired = Alias{`$EXPIRED`, FilterImplementation, nil, func([]string) string { return HasExpired().String() }}.Register()
var not_expired = Alias{`$NOT_EXPIRED`, FilterImplementation, nil, func([]string) string { return Not(HasExpired()).String() }}.Register()
var id = Alias{`$ID`, FilterImplementation, []string{"id"}, func(params []string) string { return ByID(params[0]).String() }}.Register()
var member_of = Alias{`$MEMBER_OF`, FilterImplementation, []string{"dn", "?recurse"}, func(params []string) string {
	if len(params) > 1 {
		if p, err := strconv.ParseBool(params[1]); err == nil && p {
			return MemberOf(params[0], true).String()
		}
	}

	return MemberOf(params[0], false).String()
}}.Register()

// composite filters
var and = Alias{`$AND`, FilterComposition, []string{"(filter1)", "...", "(filterN)"}, func(params []string) string { return `(&` + strings.Join(params, "") + `)` }}.Register()
var not = Alias{`$NOT`, FilterComposition, []string{"(filter)"}, func(params []string) string { return `(!` + strings.Join(params, "") + `)` }}.Register()
var or = Alias{`$OR`, FilterComposition, []string{"(filter1)", "...", "(filterN)"}, func(params []string) string { return `(|` + strings.Join(params, "") + `)` }}.Register()

// matching rules
var band = Alias{`$BAND`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_AND) }}.Register()
var bor = Alias{`$BOR`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_OR) }}.Register()
var recursive = Alias{`$RECURSIVE`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_IN_CHAIN) }}.Register()
var data = Alias{`$DATA`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA) }}.Register()

// custom filters
var attr = Alias{`$ATTR`, FilterImplementation, []string{"attribute", "value", "?operator", "?rule"}, func(params []string) string {
	attr := attributes.Lookup(params[0])
	if attr == nil {
		attr = &attributes.Attribute{LDAPDisplayName: params[0], Type: attributes.TypeRaw}
	}

	filter := &Filter{Attribute: *attr}
	if len(params) > 1 {
		filter.Value = params[1]
	}

	if len(params) > 2 {
		filter.Value = params[2] + filter.Value
	}

	if len(params) > 3 {
		filter.Rule = attributes.MatchingRule(params[3])
	}

	return filter.String()
}}.Register()

func Enabled() Alias     { return enabled }
func Disabled() Alias    { return disabled }
func Group() Alias       { return group }
func User() Alias        { return user }
func Dc() Alias          { return dc }
func Expired() Alias     { return expired }
func Not_expired() Alias { return not_expired }
func Id() Alias          { return id }
func Member_of() Alias   { return member_of }
func AndAlias() Alias    { return and }
func NotAlias() Alias    { return not }
func OrAlias() Alias     { return or }
func Band() Alias        { return band }
func Bor() Alias         { return bor }
func Recursive() Alias   { return recursive }
func Data() Alias        { return data }
func Attr() Alias        { return attr }
