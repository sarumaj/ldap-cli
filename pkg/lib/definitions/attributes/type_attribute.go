package attributes

import (
	"slices"

	libutil "github.com/sarumaj/ldap-cli/v2/pkg/lib/util"
)

// Attribute is an LDAP attribute
type Attribute struct {
	// Alias is the alias of the attribute
	Alias string
	// LDAPDisplayName is the LDAP display name of the attribute
	LDAPDisplayName string
	// PrettyName is the pretty name of the attribute (CamelCase)
	PrettyName string
	// Type is the type of the attribute defined in RFC 4517
	Type Type
}

// Parse parses a value and stores it in a map of attributes
//
//gocyclo:ignore
func (a Attribute) Parse(values []string, attrMap *Map) {
	if len(values) == 0 || attrMap == nil {
		return
	}

	if parser, ok := map[Type]func(Attribute, []string){
		TypeBool:        attrMap.ParseBool,
		TypeDecimal:     attrMap.ParseDecimal,
		TypeGroupType:   attrMap.ParseGroupType,
		TypeHexString:   func(a Attribute, s []string) { (*attrMap)[a] = libutil.Hexify(values[0]) },
		TypeInt:         attrMap.ParseInt,
		TypeIPv4Address: attrMap.ParseIPv4Address,
		TypeRaw: func(a Attribute, s []string) {
			if len(values) == 1 {
				(*attrMap)[a] = values[0]
			} else {
				(*attrMap)[a] = values
			}
		},
		TypeSAMaccountType:     attrMap.ParseSAMAccountType,
		TypeString:             func(a Attribute, s []string) { (*attrMap)[a] = values[0] },
		TypeStringSlice:        func(a Attribute, s []string) { (*attrMap)[a] = values },
		TypeTime:               attrMap.ParseTime,
		TypeUserAccountControl: attrMap.ParseUserAccountControl,
	}[a.Type]; ok {
		parser(a, values)
	}
}

// Register registers an attribute in the registry
func (a Attribute) Register() Attribute {
	registry.Append(a)
	return a
}

// String returns the string representation of an attribute (first its pretty name, then its LDAP display name)
func (a Attribute) String() string {
	if a.PrettyName != "" {
		return a.PrettyName
	}

	return a.LDAPDisplayName
}

// Attributes is a slice of attributes
type Attributes []Attribute

// Append appends attributes to the slice only if they are not already present
func (a *Attributes) Append(attrs ...Attribute) {
	for _, add := range attrs {
		seen := false
		for _, there := range *a {
			if add == there {
				seen = true
				break
			}
		}

		if !seen {
			*a = append(*a, add)
		}
	}

	a.Sort()
}

// Sort sorts the attributes by their string representation
func (a Attributes) Sort() {
	slices.SortStableFunc(a, func(a, b Attribute) int {
		l, r := a.String(), b.String()

		if l > r {
			return 1
		}

		return -1
	})
}

// ToAttributeList returns a list of attributes as strings
func (a Attributes) ToAttributeList() (list []string) {
	seen := make(map[Attribute]bool)
	for _, attr := range a {
		if seen[attr] {
			continue
		}

		if attr.LDAPDisplayName != "" {
			list = append(list, libutil.ToTitleNoLower(attr.LDAPDisplayName))
		} else {
			list = append(list, libutil.ToTitleNoLower(attr.PrettyName))
		}
		seen[attr] = true
	}

	slices.Sort(list)
	return list
}
