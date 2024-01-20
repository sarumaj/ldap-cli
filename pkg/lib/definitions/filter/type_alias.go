package filter

import (
	"bytes"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const (
	// FilterImplementation is the kind of an alias that is a custom filter implementation
	FilterImplementation Kind = "filter implementation"
	FilterComposition    Kind = "filter composition"
	// MatchingRule is the kind of an alias that is a matching rule bit mask
	MatchingRule Kind = "matching rule"
)

// validAliasCharset matches valid characters for an alias expression
var validAliasCharset = regexp.MustCompile(`^[A-Z_]+$`)

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

// findMatches finds all matches of a composite alias expression in a sequence of bytes
func (a Alias) findMatches(indexes [][]int, raw []byte) (matches [][]byte) {
	for _, index := range indexes {
		// find the end of the alias expression by searching for the closing parenthesis
		for i, l, r := index[1], 0, 0; i < len(raw); i++ {
			switch c := raw[i]; {

			case c == '(':
				l += 1

			case c == ')':
				r += 1

			}

			if l*r > 0 && l == r {
				match := raw[index[0] : i+1]
				matches = append(matches, match)
				break

			}
		}
	}
	return
}

// findOccurences returns the indexes of all occurrences of an alias expression in a sequence of bytes
func (a Alias) findOccurences(raw []byte) (indexes [][]int) {
	for _, index := range regexp.MustCompile(regexp.QuoteMeta(a.ID)).FindAllIndex(raw, -1) {
		// skip if the alias expression is part of a longer alias expression (avoid collisions of aliases)
		// raw[index[0]+1:index[1]+1] is the alias expression without the first character
		if index[1] < len(raw) && validAliasCharset.Match(raw[index[0]+1:index[1]+1]) {
			continue
		}

		indexes = append(indexes, index)
	}

	return
}

// findSplitPositions finds the split positions of the parameters of a composite alias expression
func (a Alias) findSplitPositions(match []byte) (splits []int) {
	// strip the alias ID from the parameters
	// parameters are supposed to be wrapped in parenthesis and separated by semicolons
	params := bytes.TrimPrefix(match, []byte(a.ID))

	// split the parameters by semicolons
	// the splits can be nested in parenthesis
	// the difficulty is to find the correct split positions
	for i, l, r := 0, 0, 0; i < len(params); i++ {
		switch c := params[i]; {

		case c == '(':
			l += 1

		case c == ')':
			r += 1

		case c == ';' && l-r == 1:
			// when the number of opening parenthesis is greater than the number of the closing parenthesis by one,
			// the semicolon is a split for given alias expression (lower or upper boundary)
			splits = append(splits, i)

		}

		// when the number of opening parenthesis is equal to the number of the closing parenthesis,
		// the last split position (upper boundary) is the end of the alias expression
		if l*r > 0 && l == r {
			splits = append(splits, i+1)
			break
		}

	}

	return
}

// isComposite returns true if the alias is a composite alias
func (a Alias) isComposite() bool { return len(a.Parameters) > 0 }

// Register registers an alias in the registry
func (a Alias) Register() Alias {
	registry = append(registry, a)
	return a
}

// replace replaces all occurrences of an alias expression in a sequence of bytes
func (a Alias) replace(raw []byte) []byte {
	indexes := a.findOccurences(raw)
	replaced := len(indexes) > 0

	switch {

	case a.isComposite():
		// composite alias substitution
		matches := a.findMatches(indexes, raw)

		// substitute all occurrences of the alias expression
		for _, match := range matches {
			// strip the alias ID from the parameters
			// parameters are supposed to be wrapped in parenthesis and separated by semicolons
			params := bytes.TrimPrefix(match, []byte(a.ID))

			// split the parameters by semicolons
			// the splits can be nested in parenthesis
			// the difficulty is to find the correct split positions
			splits := a.findSplitPositions(match)

			// split the parameters at split positions
			params_list := a.splitParameters(params, splits)

			raw = bytes.ReplaceAll(raw, match, []byte(a.Substitution(params_list)))
		}

		replaced = replaced && len(matches) > 0

	default:
		// simple alias substitution
		var span int // span is the number of bytes added or removed during substitution
		for _, index := range indexes {
			// substitute the alias expression
			substitution := []byte(a.Substitution(nil))
			tail := append(substitution, raw[index[1]-span:]...)
			raw = append(raw[:index[0]-span], tail...)

			// update span to account for the substitution
			span += index[1] - index[0] - len(substitution)
		}

	}

	// avoid infinite recursion by checking if the alias expression has been even replaced
	if !replaced {
		return raw
	}

	// check if there are any aliases left to substitute and call recursively
	if left := a.findOccurences(raw); len(left) > 0 {
		return a.replace(raw)
	}

	return raw
}

// splitParameters splits the parameters of an alias expression
func (a Alias) splitParameters(params []byte, splits []int) (params_list []string) {
	for i, start := 0, 0; i < len(splits); i, start = i+1, splits[i]+1 {
		param := params[start:splits[i]]

		// trim opening parenthesis from the first parameter
		if i == 0 {
			param = bytes.TrimPrefix(param, []byte{'('})
		}

		// trim closing parenthesis from the last parameter
		if i == len(splits)-1 {
			param = bytes.TrimSuffix(param, []byte{')'})
		}

		// strip the parameter from the leading and trailing whitespace and
		// append it to the parameter list
		params_list = append(params_list, string(bytes.TrimSpace(param)))
	}

	return
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
	list := make([]Alias, len(registry))
	_ = copy(list, registry)

	slices.SortStableFunc(list, func(a, b Alias) int {
		if string(a.Kind) == string(b.Kind) && a.ID > b.ID {
			return 1
		}

		return -1
	})

	return list
}

// ReplaceAliases replaces alias expressions in the input string
func ReplaceAliases(in string) string {
	raw_bytes := []byte(in)
	for _, alias := range registry {
		raw_bytes = alias.replace(raw_bytes)
	}

	return string(raw_bytes)
}
