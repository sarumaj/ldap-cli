package attributes

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"

	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type Map map[Attribute]any

func (attrMap Map) Keys() (keys Attributes) {
	for a := range attrMap {
		keys.Append(a)
	}

	keys.Sort()
	return
}

func (attrMap *Map) ParseBool(a Attribute, values []string) {
	if parsed, err := strconv.ParseBool(values[0]); err == nil {
		(*attrMap)[a] = parsed
	} else {
		(*attrMap)[a] = values
	}

}

func (attrMap *Map) ParseDecimal(a Attribute, values []string) {
	if parsed, err := strconv.ParseFloat(values[0], 64); err == nil {
		(*attrMap)[a] = parsed
	} else {
		(*attrMap)[a] = values
	}
}

func (attrMap *Map) ParseGroupType(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		(*attrMap)[a] = FlagsetGroupType(parsed).Eval()
	} else {
		(*attrMap)[a] = values
	}
}

func (attrMap *Map) ParseInt(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		(*attrMap)[a] = parsed
	} else {
		(*attrMap)[a] = values
	}
}

func (attrMap *Map) ParseIPv4Address(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(parsed))
		(*attrMap)[a] = ip
	} else {
		(*attrMap)[a] = values
	}
}

func (attrMap *Map) ParseTime(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(strings.Split(values[0], ".")[0], 10, 64); err == nil {
		(*attrMap)[a] = libutil.TimeAfter1601(parsed)
	} else {
		(*attrMap)[a] = values
	}
}

func (attrMap *Map) ParseSAMAccountType(a Attribute, values []string) {
	if parsed, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		(*attrMap)[a] = FlagSAMAccountType(parsed).Eval()
	} else {
		(*attrMap)[a] = values
	}
}

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

type Maps []map[Attribute]any
