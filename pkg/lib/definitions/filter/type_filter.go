package filter

import (
	"fmt"
	"regexp"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

// complexFilterSyntax is used to identify complex filters (internal agreement)
const complexFilterSyntax = "complex"

// notDefaultOperatorWithValueRegex matches operators that are not the default operator ("=") for the given attribute type
var notDefaultOperatorWithValueRegex = regexp.MustCompile(`^(?P<Operator>[~<>]=)` + `(?P<Value>.*)$`)

// Filter is used to define an LDAP filter
type Filter struct {
	// Attribute is the attribute to filter on
	Attribute attributes.Attribute
	// Value is the value to filter for
	Value string
	// Rule is the matching rule bit mask to use
	Rule attributes.MatchingRule
}

// ExpandAlias expands an filter to a filter by itself or its attribute alias
func (o Filter) ExpandAlias() Filter {
	if o.Attribute.Alias != "" {
		return Or(o, Filter{
			Attribute: attributes.Attribute{LDAPDisplayName: o.Attribute.Alias, Type: o.Attribute.Type},
			Value:     o.Value,
			Rule:      o.Rule,
		})
	}

	return o
}

// String returns a string representation of a filter (LDAP filter syntax)
func (o Filter) String() string {
	if o.Attribute.LDAPDisplayName == complexFilterSyntax {
		return o.Value
	}

	op, value := "=", strings.TrimPrefix(o.Value, "=")
	if notDefaultOperatorWithValueRegex.MatchString(value) {
		op, value = value[:2], value[2:]
	}

	switch o.Attribute.Type {

	case // gt, lt, and proximity operators not allowed
		attributes.TypeBool:

		op, value = strings.TrimLeft(op, "<>~"), strings.ToUpper(value)

	case // gt or lt operators not allowed
		attributes.TypeHexString,
		attributes.TypeString,
		attributes.TypeStringSlice:

		op = strings.TrimLeft(op, "<>")

	case // proximity operator not allowed
		attributes.TypeDecimal,
		attributes.TypeGroupType,
		attributes.TypeIPv4Address,
		attributes.TypeInt,
		attributes.TypeSAMaccountType,
		attributes.TypeTime,
		attributes.TypeUserAccountControl:

		op = strings.TrimLeft(op, "~")

	}

	if o.Attribute.Type == attributes.TypeBool {
		value = strings.ToUpper(value)
	}

	switch {

	case o.Rule != "":
		return fmt.Sprintf("(%s:%s:%s%s)", libutil.ToTitleNoLower(o.Attribute.LDAPDisplayName), o.Rule, op, value)

	default:
		return fmt.Sprintf("(%s%s%s)", libutil.ToTitleNoLower(o.Attribute.LDAPDisplayName), op, value)

	}
}

// And returns a filter that matches all given filters
func And(property Filter, properties ...Filter) Filter {
	return complexFilter('&', property, properties...)
}

// complexFilter is used to build complex filters
func complexFilter(operator rune, property Filter, properties ...Filter) Filter {
	if len(properties) == 0 {
		return property
	}

	var values []string
	seen := make(map[string]bool)
	for _, property := range append([]Filter{property}, properties...) {

		if v, ok := seen[property.String()]; ok && v {
			continue
		}

		key := property.String()
		values, seen[key] = append(values, key), true
	}

	if len(values) == 1 {
		return property
	}

	return Filter{
		Attribute: attributes.Attribute{LDAPDisplayName: complexFilterSyntax},
		Value:     "(" + string(operator) + strings.Join(values, "") + ")",
	}
}

// EscapeFilter returns a string representation of a filter with escaped values according to RFC 4515
func EscapeFilter(filter string) string { return ldap.EscapeFilter(filter) }

// Not returns a filter that matches the opposite of the given filter
func Not(property Filter) Filter {
	if true &&
		property.Attribute.LDAPDisplayName == complexFilterSyntax &&
		strings.HasPrefix(property.Value, "(!") &&
		strings.HasSuffix(property.Value, ")") {

		return Filter{
			Attribute: property.Attribute,
			Value:     property.Value[2 : len(property.Value)-1],
		}
	}

	return Filter{
		Attribute: attributes.Attribute{LDAPDisplayName: complexFilterSyntax},
		Value:     "(!" + property.String() + ")",
	}
}

// Or returns a filter that matches any of the given filters
func Or(property Filter, properties ...Filter) Filter {
	return complexFilter('|', property, properties...)
}
