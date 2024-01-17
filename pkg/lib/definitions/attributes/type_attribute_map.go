package attributes

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"

	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

// Map is a map of attributes to values
type Map map[Attribute]any

// Keys returns the keys of a map of attributes
func (attrMap Map) Keys() (keys Attributes) {
	for a := range attrMap {
		keys.Append(a)
	}

	keys.Sort()
	return
}

// ParseBool parses a boolean value
func (attrMap *Map) ParseBool(a Attribute, values []string) {
	if parsed, err := strconv.ParseBool(values[0]); err == nil {
		(*attrMap)[a] = parsed
	} else {
		(*attrMap)[a] = values
	}

}

// ParseDecimal parses a decimal value
func (attrMap *Map) ParseDecimal(a Attribute, values []string) {
	if parsed, err := strconv.ParseFloat(values[0], 64); err == nil {
		(*attrMap)[a] = parsed
	} else {
		(*attrMap)[a] = values
	}
}

// ParseGroupType parses a group type value
func (attrMap *Map) ParseGroupType(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		(*attrMap)[a] = FlagsetGroupType(parsed).Eval()
	} else {
		(*attrMap)[a] = values
	}
}

// ParseInt parses an integer value
func (attrMap *Map) ParseInt(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		(*attrMap)[a] = parsed
	} else {
		(*attrMap)[a] = values
	}
}

// ParseIPv4Address parses an IPv4 address value
func (attrMap *Map) ParseIPv4Address(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(parsed))
		(*attrMap)[a] = ip
	} else {
		(*attrMap)[a] = values
	}
}

// ParseTime parses a time value
func (attrMap *Map) ParseTime(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(strings.Split(values[0], ".")[0], 10, 64); err == nil {
		(*attrMap)[a] = libutil.TimeAfter1601(parsed)
	} else {
		(*attrMap)[a] = values
	}
}

// ParseSAMAccountName parses a SAM account name (SAN) value
func (attrMap *Map) ParseSAMAccountType(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		(*attrMap)[a] = FlagSAMAccountType(parsed).Eval()
	} else {
		(*attrMap)[a] = values
	}
}

// ParseUserAccountControl parses a user account control (AUC) value
func (attrMap *Map) ParseUserAccountControl(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		userAccountControl := FlagsetUserAccountControl(parsed)
		(*attrMap)[a] = userAccountControl.Eval()
		(*attrMap)[Raw("", "Enabled", TypeBool)] = userAccountControl&USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE == 0
		(*attrMap)[Raw("", "LockedOut", TypeBool)] = userAccountControl&USER_ACCOUNT_CONTROL_LOCKOUT != 0
	} else {
		(*attrMap)[a] = values
	}
}

// Maps is a slice of maps of attributes to values
type Maps []Map
