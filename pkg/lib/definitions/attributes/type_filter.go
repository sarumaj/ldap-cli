package attributes

import (
	"fmt"
	"regexp"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

const complexFilterSyntax = "complex"

var nonDefaultOperatorWithValueRegex = regexp.MustCompile(`^(?P<Operator>[~<>]=)(?P<Value>.*)$`)

type Filter struct {
	Attribute Attribute
	Value     string
	Rule      MatchingRule
}

func (o Filter) String() string {
	op, value := "=", o.Value
	if nonDefaultOperatorWithValueRegex.MatchString(value) {
		op, value = value[:2], value[2:]
	}

	switch {

	case o.Attribute == complexFilterSyntax:
		return o.Value

	case o.Rule != "":
		return fmt.Sprintf("(%s:%s:%s%s)", o.Attribute, o.Rule, op, value)

	default:
		return fmt.Sprintf("(%s%s%s)", o.Attribute, op, value)

	}
}

// Build complex filter from filters, where all must match
func And(property Filter, properties ...Filter) Filter {
	var value string
	for _, property := range append([]Filter{property}, properties...) {
		value += property.String()
	}

	return Filter{
		Attribute: complexFilterSyntax,
		Value:     "(&" + value + ")",
	}
}

// Escape special characters as specified in RFC4515
func EscapeFilter(filter string) string { return ldap.EscapeFilter(filter) }

func HasNotExpired(strict bool) Filter {
	filter := Or(
		Filter{Attribute: AttributeAccountExpires, Value: fmt.Sprint(0)},
		Filter{Attribute: AttributeAccountExpires, Value: fmt.Sprint(1<<63 - 1)},
		Filter{
			Attribute: AttributeAccountExpires,
			Value:     fmt.Sprintf(">=%d", util.TimeSince1601().Nanoseconds()/100),
		},
	)

	if strict {
		return And(filter, Filter{Attribute: AttributeAccountExpires, Value: "*"})
	}

	return Or(filter, Not(Filter{Attribute: AttributeAccountExpires, Value: "*"}))
}

func IsDomainController() Filter {
	return And(
		Filter{AttributeObjectClass, "computer", ""},
		Filter{AttributeUserAccountControl, fmt.Sprintf("%d", USER_ACCOUNT_CONTROL_SERVER_TRUST_ACCOUNT), LDAP_MATCHING_RULE_BIT_AND},
	)
}

func IsEnabled() Filter {
	return Not(Filter{AttributeUserAccountControl, "2", LDAP_MATCHING_RULE_BIT_AND})
}

func IsGroup() Filter { return Filter{AttributeObjectClass, "group", ""} }

func IsUser() Filter { return Filter{AttributeObjectClass, "user", ""} }

// Negate filter
func Not(property Filter) Filter {
	return Filter{
		Attribute: complexFilterSyntax,
		Value:     "(!" + property.String() + ")",
	}
}

// Build complex filter from filters, where at least one must match
func Or(property Filter, properties ...Filter) Filter {
	var value string
	for _, property := range append([]Filter{property}, properties...) {
		value += property.String()
	}

	return Filter{
		Attribute: complexFilterSyntax,
		Value:     "(|" + value + ")",
	}
}
