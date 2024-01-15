package filter

import (
	"fmt"
	"regexp"
	"strings"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

var validSimpleFilterRegex = regexp.MustCompile(`^\((?P<Attribute>[\w\-]+)(?::(?P<Rule>(?:\d+\.){6}\d+):)?(?P<Operator>[~<>]?=)(?P<Value>.*)\)$`)
var validComplexFilterRegex = regexp.MustCompile(`^\((?P<Logic>[!&\|])\((?P<Filters>.+)\)\)$`)

func balanceClosingParenthesis(in string) string {
	// count continuous consecutive closing parenthesis to determine complexity level of the filter
	var highestCount, currentCount, leftCount, rightCount int
	var notPadded bool
	for i, c := range in {
		switch {
		case i == 0 && c == '(': // compensate for first opening parenthesis
			leftCount, rightCount = -1, -1

		case i == 0 && c != '(': // compensate for first non-opening parenthesis
			notPadded = true

		case i == len(in)-1 && c != ')': // compensate for last non-closing parenthesis
			notPadded = notPadded && true

		case c == ')': // update counters and save highest count
			rightCount, currentCount = rightCount+1, currentCount+1
			if currentCount > highestCount {
				highestCount = currentCount
			}

		case c == '(': // update counter
			leftCount += 1
			fallthrough

		default:
			currentCount = 0
		}
	}

	var padLeft, padRight string
	switch balance := highestCount + leftCount - rightCount; {
	case balance > 0 && notPadded: // only if the input is not padded and closing parenthesis are unbalanced
		padLeft, padRight = "(", strings.Repeat(")", balance)

	case notPadded: // handle case when input was not padded
		padLeft, padRight = "(", ")"

	}

	return padLeft + in + padRight
}

func ParseRaw(raw string) (*Filter, error) {
	switch raw = ReplaceAliases(raw); {

	case validSimpleFilterRegex.MatchString(raw):
		matches := validSimpleFilterRegex.FindStringSubmatch(raw)[1:]
		name, rule, value := matches[0], matches[1], matches[2]+matches[3]

		attr := attributes.Lookup(name)

		// not found, so a custom attribute has been used
		if attr == nil {
			attr = &attributes.Attribute{LDAPDisplayName: name, Type: attributes.TypeRaw}
		}

		// aliased attribute, enforce not-aliased rendering
		if attr.Alias == name {
			attr = &attributes.Attribute{LDAPDisplayName: name, Type: attr.Type}
		}

		return &Filter{
			Attribute: *attr,
			Rule:      attributes.MatchingRule(rule),
			Value:     value,
		}, nil

	case validComplexFilterRegex.MatchString(raw):
		matches := validComplexFilterRegex.FindStringSubmatch(raw)[1:]
		junction, content := matches[0], matches[1]

		// since content is a complex filter stripped of parenthesis, restore padding and bring parenthesis in balance
		padded := balanceClosingParenthesis(content)

		// divide complex filter into simpler filters
		components := splitFilterComponents(padded)

		// parse simplified filters recursively
		var filters []Filter
		for _, component := range components {
			filter, err := ParseRaw(component)
			if err != nil {
				return nil, err
			}

			filters = append(filters, *filter)
		}

		// reassembly complex filter
		if complexFn, ok := map[rune]func() *Filter{
			'&': func() *Filter { f := And(filters[0], filters[1:]...); return &f },
			'|': func() *Filter { f := Or(filters[0], filters[1:]...); return &f },
			'!': func() *Filter { f := Not(filters[0]); return &f },
		}[rune(junction[0])]; ok && len(filters) > 0 {

			return complexFn(), nil
		}
	}

	return nil, fmt.Errorf("%w: %s", libutil.ErrInvalidFilter, raw)
}

func splitFilterComponents(in string) []string {
	// determine split indexes
	var openings, endings int
	var indexes []int
	for i, c := range in {
		switch c {
		case '(':
			openings += 1

		case ')':
			endings += 1

		}

		// split should appear whenever openings and endings are in balance
		if openings*endings > 0 && endings == openings {
			indexes = append(indexes, i)
			openings, endings = 0, 0
		}
	}

	// perform splits
	var parts []string
	for i, index := range indexes {
		if i == 0 {
			parts = append(parts, in[:index+1])
		} else {
			parts = append(parts, in[indexes[i-1]+1:index+1])
		}

	}

	return parts
}
