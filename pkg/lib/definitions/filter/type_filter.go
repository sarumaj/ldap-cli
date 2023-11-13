package filter

import (
	"fmt"
	"regexp"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

const complexFilterSyntax = "complex"

var notDefaultOperatorWithValueRegex = regexp.MustCompile(`^(?P<Operator>[~<>]=)` + `(?P<Value>.*)$`)

type Filter struct {
	Attribute attributes.Attribute
	Value     string
	Rule      attributes.MatchingRule
}

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

// Build complex filter from filters, where all must match
func And(property Filter, properties ...Filter) Filter {
	return complexFilter('&', property, properties...)
}

func complexFilter(operator rune, property Filter, properties ...Filter) Filter {
	if len(properties) == 0 {
		return property
	}

	var value string
	for _, property := range append([]Filter{property}, properties...) {
		value += property.String()
	}

	return Filter{
		Attribute: attributes.Attribute{LDAPDisplayName: complexFilterSyntax},
		Value:     "(" + string(operator) + value + ")",
	}
}

// Escape special characters as specified in RFC4515
func EscapeFilter(filter string) string { return ldap.EscapeFilter(filter) }

// Negate filter
func Not(property Filter) Filter {
	return Filter{
		Attribute: attributes.Attribute{LDAPDisplayName: complexFilterSyntax},
		Value:     "(!" + property.String() + ")",
	}
}

// Build complex filter from filters, where at least one must match
func Or(property Filter, properties ...Filter) Filter {
	return complexFilter('|', property, properties...)
}
