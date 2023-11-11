package attributes

import (
	"encoding/binary"
	"net"
	"slices"
	"strconv"
	"strings"

	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type Attribute struct {
	Alias           string
	LDAPDisplayName string
	PrettyName      string
	Type            Type
}

func (a Attribute) Parse(values []string, attrMap *Map) {
	if len(values) == 0 || attrMap == nil {
		return
	}

	switch a.Type {
	case TypeBool:
		parsed, err := strconv.ParseBool(values[0])
		if err == nil {
			(*attrMap)[a] = parsed
		} else {
			(*attrMap)[a] = values
		}

	case TypeDecimal:
		parsed, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			(*attrMap)[a] = parsed
		} else {
			(*attrMap)[a] = values
		}

	case TypeGroupType:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = FlagsetGroupType(parsed).Eval()
		} else {
			(*attrMap)[a] = values
		}

	case TypeHexString:
		(*attrMap)[a] = libutil.Hexify(values[0])

	case TypeInt:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = parsed
		} else {
			(*attrMap)[a] = values
		}

	case TypeIPv4Address:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, uint32(parsed))
			(*attrMap)[a] = ip
		} else {
			(*attrMap)[a] = values
		}

	case TypeRaw:
		if len(values) == 1 {
			(*attrMap)[a] = values[0]
		} else {
			(*attrMap)[a] = values
		}

	case TypeSAMaccountType:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = FlagSAMAccountType(parsed).Eval()
		} else {
			(*attrMap)[a] = values
		}

	case TypeString:
		(*attrMap)[a] = values[0]

	case TypeStringSlice:
		(*attrMap)[a] = values

	case TypeTime:
		parsed, err := strconv.ParseInt(strings.Split(values[0], ".")[0], 10, 64)
		if err == nil {
			(*attrMap)[a] = libutil.TimeAfter1601(parsed)
		} else {
			(*attrMap)[a] = values
		}

	case TypeUserAccountControl:
		parsed, err := strconv.ParseInt(values[0], 10, 64)
		if err == nil {
			userAccountControl := FlagsetUserAccountControl(parsed)
			(*attrMap)[a] = userAccountControl.Eval()
			(*attrMap)[Raw("", "Enabled", TypeBool)] = userAccountControl&USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE == 0
			(*attrMap)[Raw("", "LockedOut", TypeBool)] = userAccountControl&USER_ACCOUNT_CONTROL_LOCKOUT != 0
		} else {
			(*attrMap)[a] = values
		}

	default:
		return

	}
}

func (a Attribute) Register() Attribute {
	registry.Append(a)
	return a
}

func (a Attribute) String() string {
	if a.PrettyName != "" {
		return a.PrettyName
	}

	return a.LDAPDisplayName
}

type Attributes []Attribute

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

func (a Attributes) Sort() {
	slices.SortStableFunc(a, func(a, b Attribute) int {
		l, r := a.String(), b.String()

		if l > r {
			return 1
		}

		return -1
	})
}

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
