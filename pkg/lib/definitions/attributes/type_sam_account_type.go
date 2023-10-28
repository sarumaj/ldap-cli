package attributes

import "strings"

// https://docs.microsoft.com/en-us/windows/win32/adschema/a-samaccounttype
const (
	SAM_ACCOUNT_TYPE_DOMAIN_OBJECT             FlagSAMAccountType = 0x00000000
	SAM_ACCOUNT_TYPE_GROUP_OBJECT              FlagSAMAccountType = 0x10000000
	SAM_ACCOUNT_TYPE_NON_SECURITY_GROUP_OBJECT FlagSAMAccountType = 0x10000001
	SAM_ACCOUNT_TYPE_ALIAS_OBJECT              FlagSAMAccountType = 0x20000000
	SAM_ACCOUNT_TYPE_NON_SECURITY_ALIAS_OBJECT FlagSAMAccountType = 0x20000001
	SAM_ACCOUNT_TYPE_USER_OBJECT               FlagSAMAccountType = 0x30000000
	SAM_ACCOUNT_TYPE_NORMAL_USER_ACCOUNT       FlagSAMAccountType = 0x30000000
	SAM_ACCOUNT_TYPE_MACHINE_ACCOUNT           FlagSAMAccountType = 0x30000001
	SAM_ACCOUNT_TYPE_TRUST_ACCOUNT             FlagSAMAccountType = 0x30000002
	SAM_ACCOUNT_TYPE_APP_BASIC_GROUP           FlagSAMAccountType = 0x40000000
	SAM_ACCOUNT_TYPE_APP_QUERY_GROUP           FlagSAMAccountType = 0x40000001
	SAM_ACCOUNT_TYPE_ACCOUNT_TYPE_MAX          FlagSAMAccountType = 0x7FFFFFFF
)

var samAccountTypeToString = map[FlagSAMAccountType][]string{
	SAM_ACCOUNT_TYPE_DOMAIN_OBJECT:             {"DOMAIN_OBJECT"},
	SAM_ACCOUNT_TYPE_GROUP_OBJECT:              {"GROUP_OBJECT"},
	SAM_ACCOUNT_TYPE_NON_SECURITY_GROUP_OBJECT: {"NON_SECURITY_GROUP_OBJECT"},
	SAM_ACCOUNT_TYPE_ALIAS_OBJECT:              {"ALIAS_OBJECT"},
	SAM_ACCOUNT_TYPE_NON_SECURITY_ALIAS_OBJECT: {"NON_SECURITY_ALIAS_OBJECT"},
	SAM_ACCOUNT_TYPE_USER_OBJECT:               {"NORMAL_USER_ACCOUNT", "USER_OBJECT"},
	SAM_ACCOUNT_TYPE_MACHINE_ACCOUNT:           {"MACHINE_ACCOUNT"},
	SAM_ACCOUNT_TYPE_TRUST_ACCOUNT:             {"TRUST_ACCOUNT"},
	SAM_ACCOUNT_TYPE_APP_BASIC_GROUP:           {"APP_BASIC_GROUP"},
	SAM_ACCOUNT_TYPE_APP_QUERY_GROUP:           {"APP_QUERY_GROUP"},
	SAM_ACCOUNT_TYPE_ACCOUNT_TYPE_MAX:          {"ACCOUNT_TYPE_MAX"},
}

type FlagSAMAccountType uint32

func (v FlagSAMAccountType) Eval() []string {
	for key, value := range samAccountTypeToString {
		if v == key {
			return value
		}
	}

	return nil
}

func (s FlagSAMAccountType) String() string { return strings.Join(s.Strings(), " | ") }

func (s FlagSAMAccountType) Strings() []string {
	if v, ok := samAccountTypeToString[s]; ok {
		return v
	}

	return nil
}
