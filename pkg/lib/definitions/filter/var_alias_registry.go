package filter

import (
	"strconv"
	"strings"

	attributes "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/attributes"
)

var (
	// filters
	enabled     = Alias{`$ENABLED`, FilterImplementation, nil, func([]string) string { return IsEnabled().String() }}.Register()
	disabled    = Alias{`$DISABLED`, FilterImplementation, nil, func([]string) string { return Not(IsEnabled()).String() }}.Register()
	group       = Alias{`$GROUP`, FilterImplementation, nil, func([]string) string { return IsGroup().String() }}.Register()
	user        = Alias{`$USER`, FilterImplementation, nil, func([]string) string { return IsUser().String() }}.Register()
	dc          = Alias{`$DC`, FilterImplementation, nil, func([]string) string { return IsDomainController().String() }}.Register()
	expired     = Alias{`$EXPIRED`, FilterImplementation, nil, func([]string) string { return HasExpired().String() }}.Register()
	not_expired = Alias{`$NOT_EXPIRED`, FilterImplementation, nil, func([]string) string { return Not(HasExpired()).String() }}.Register()
	id          = Alias{`$ID`, FilterImplementation, []string{"<id>"}, func(params []string) string { return ByID(params[0]).String() }}.Register()
	member_of   = Alias{`$MEMBER_OF`, FilterImplementation, []string{"<dn>", "<?recurse>"}, func(params []string) string {
		if len(params) > 1 {
			if p, err := strconv.ParseBool(params[1]); err == nil && p {
				return MemberOf(params[0], true).String()
			}
		}

		return MemberOf(params[0], false).String()
	}}.Register()

	// composite filters
	and = Alias{`$AND`, FilterComposition, []string{"(filter1)", "...", "(filterN)"}, func(params []string) string { return `(&` + strings.Join(params, "") + `)` }}.Register()
	not = Alias{`$NOT`, FilterComposition, []string{"(filter)"}, func(params []string) string { return `(!` + strings.Join(params, "") + `)` }}.Register()
	or  = Alias{`$OR`, FilterComposition, []string{"(filter1)", "...", "(filterN)"}, func(params []string) string { return `(|` + strings.Join(params, "") + `)` }}.Register()

	// matching rules
	band      = Alias{`$BAND`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_AND) }}.Register()
	bor       = Alias{`$BOR`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_BIT_OR) }}.Register()
	recursive = Alias{`$RECURSIVE`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_IN_CHAIN) }}.Register()
	data      = Alias{`$DATA`, MatchingRule, nil, func([]string) string { return string(attributes.LDAP_MATCHING_RULE_DN_WITH_DATA) }}.Register()

	// custom filters
	attr = Alias{`$ATTR`, FilterImplementation, []string{"<attribute>", "<value>", "<?operator>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 4)
		return customAliasSubstitution(params[0], params[1], params[2], params[3])
	}}.Register()

	equals = Alias{`$EQ`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], params[1], "=", params[2])
	}}.Register()

	like = Alias{`$LIKE`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], params[1], "~=", params[2])
	}}.Register()

	contains = Alias{`$CONTAINS`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], "*"+strings.Trim(params[1], "*")+"*", "=", params[2])
	}}.Register()

	starts_with = Alias{`$STARTS_WITH`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], strings.TrimSuffix(params[1], "*")+"*", "=", params[2])
	}}.Register()

	ends_with = Alias{`$ENDS_WITH`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], "*"+strings.TrimPrefix(params[1], "*"), "=", params[2])
	}}.Register()

	greater_than = Alias{`$GT`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], params[1], ">", params[2])
	}}.Register()

	less_than = Alias{`$LT`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], params[1], "<", params[2])
	}}.Register()

	greater_than_or_equal = Alias{`$GTE`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], params[1], ">=", params[2])
	}}.Register()

	less_than_or_equal = Alias{`$LTE`, FilterImplementation, []string{"<attribute>", "<value>", "<?rule>"}, func(params []string) string {
		params = extendParams(params, 3)
		return customAliasSubstitution(params[0], params[1], "<=", params[2])
	}}.Register()

	exists = Alias{`$EXISTS`, FilterImplementation, []string{"<attribute>"}, func(params []string) string {
		params = extendParams(params, 1)
		return customAliasSubstitution(params[0], "*", "=", "")
	}}.Register()

	not_exists = Alias{`$NOT_EXISTS`, FilterImplementation, []string{"<attribute>"}, func(params []string) string {
		params = extendParams(params, 1)
		return "(!" + customAliasSubstitution(params[0], "*", "=", "") + ")"
	}}.Register()
)

// list of registry
var registry []Alias

func AliasForEnabled() Alias            { return enabled }
func AliasForDisabled() Alias           { return disabled }
func AliasForGroup() Alias              { return group }
func AliasForUser() Alias               { return user }
func AliasForDc() Alias                 { return dc }
func AliasForExpired() Alias            { return expired }
func AliasForNotExpired() Alias         { return not_expired }
func AliasForId() Alias                 { return id }
func AliasForMemberOf() Alias           { return member_of }
func AliasForAnd() Alias                { return and }
func AliasForNot() Alias                { return not }
func AliasForOr() Alias                 { return or }
func AliasForBand() Alias               { return band }
func AliasForBor() Alias                { return bor }
func AliasForRecursive() Alias          { return recursive }
func AliasForData() Alias               { return data }
func AliasForAttr() Alias               { return attr }
func AliasForEquals() Alias             { return equals }
func AliasForLike() Alias               { return like }
func AliasForContains() Alias           { return contains }
func AliasForStartsWith() Alias         { return starts_with }
func AliasForEndsWith() Alias           { return ends_with }
func AliasForGreaterThan() Alias        { return greater_than }
func AliasForLessThan() Alias           { return less_than }
func AliasForGreaterThanOrEqual() Alias { return greater_than_or_equal }
func AliasForLessThanOrEqual() Alias    { return less_than_or_equal }
func AliasForExists() Alias             { return exists }
func AliasForNotExists() Alias          { return not_exists }

// customAliasSubstitution returns a filter that matches the given attribute with the given value
func customAliasSubstitution(attributeName, value, operator, rule string) string {
	attr := attributes.Lookup(attributeName)
	if attr == nil {
		attr = &attributes.Attribute{LDAPDisplayName: attributeName, Type: attributes.TypeRaw}
	}

	return Filter{
		Attribute: *attr,
		Value:     operator + value,
		Rule:      attributes.MatchingRule(rule),
	}.String()
}

// extendParams extends the parameter list to the given length
func extendParams(params []string, n int) []string {
	for len(params) < n {
		params = append(params, "")
	}

	return params
}
