package filter

import (
	"fmt"
	"regexp"
	"strings"

	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

var validSimpleFilterRegex = regexp.MustCompile(`^\(` + `(?P<Attribute>[\w\-]+)` + `(?::(?P<Rule>(?:\d+\.){6}\d+):)?` + `(?P<Operator>[~<>]?=)` + `(?P<Value>.*)` + `\)$`)
var validComplexFilterRegex = regexp.MustCompile(`^\(` + `(?P<Logic>[!&\|])` + `\(` + `(?P<Filters>.+)` + `\)` + `\)$`)

func ParseRaw(raw string) (*Filter, error) {
	for _, alias := range aliases {
		raw = strings.ReplaceAll(raw, alias.Alias, alias.Substitution)
	}

	switch {

	case validSimpleFilterRegex.MatchString(raw):
		matches := validSimpleFilterRegex.FindStringSubmatch(raw)[1:]
		name, rule, value := matches[0], matches[1], matches[2]+matches[3]

		attr := attributes.Lookup(name)
		if attr == nil {
			attr = &attributes.Attribute{LDAPDisplayName: name, Type: attributes.TypeRaw}
		}

		return &Filter{
			Attribute: *attr,
			Rule:      attributes.MatchingRule(rule),
			Value:     value,
		}, nil

	case validComplexFilterRegex.MatchString(raw):
		matches := validComplexFilterRegex.FindStringSubmatch(raw)[1:]
		junction, content := matches[0], matches[1]

		var filters []Filter
		for _, raw := range strings.Split(content, ")(") {
			filter, err := ParseRaw("(" + raw + ")")
			if err != nil {
				return nil, err
			}

			filters = append(filters, *filter)
		}

		if complexFn, ok := map[rune]func() Filter{
			'&': func() Filter { return And(filters[0], filters[1:]...) },
			'|': func() Filter { return Or(filters[0], filters[1:]...) },
			'!': func() Filter { return Not(filters[0]) },
		}[rune(junction[0])]; ok && len(filters) > 0 {

			complex := complexFn()
			return &complex, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", libutil.ErrInvalidFilter, raw)
}
