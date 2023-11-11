package attributes

// https://docs.microsoft.com/en-us/windows/win32/adschema/a-samaccounttype
const (
	SAM_ACCOUNT_TYPE_DOMAIN_OBJECT             SamAccountType = 0x00000000
	SAM_ACCOUNT_TYPE_GROUP_OBJECT              SamAccountType = 0x10000000
	SAM_ACCOUNT_TYPE_NON_SECURITY_GROUP_OBJECT SamAccountType = 0x10000001
	SAM_ACCOUNT_TYPE_ALIAS_OBJECT              SamAccountType = 0x20000000
	SAM_ACCOUNT_TYPE_NON_SECURITY_ALIAS_OBJECT SamAccountType = 0x20000001
	SAM_ACCOUNT_TYPE_USER_OBJECT               SamAccountType = 0x30000000
	SAM_ACCOUNT_TYPE_NORMAL_USER_ACCOUNT       SamAccountType = 0x30000000
	SAM_ACCOUNT_TYPE_MACHINE_ACCOUNT           SamAccountType = 0x30000001
	SAM_ACCOUNT_TYPE_TRUST_ACCOUNT             SamAccountType = 0x30000002
	SAM_ACCOUNT_TYPE_APP_BASIC_GROUP           SamAccountType = 0x40000000
	SAM_ACCOUNT_TYPE_APP_QUERY_GROUP           SamAccountType = 0x40000001
	SAM_ACCOUNT_TYPE_ACCOUNT_TYPE_MAX          SamAccountType = 0x7FFFFFFF
	SAM_ACCOUNT_TYPE_UNKNOWN                   SamAccountType = 0xFFFFFFFF
)

var samAccountTypeToString = map[SamAccountType][]string{
	SAM_ACCOUNT_TYPE_DOMAIN_OBJECT:             {"DOMAIN_OBJECT"},
	SAM_ACCOUNT_TYPE_GROUP_OBJECT:              {"GROUP_OBJECT"},
	SAM_ACCOUNT_TYPE_NON_SECURITY_GROUP_OBJECT: {"NON_SECURITY_GROUP_OBJECT"},
	SAM_ACCOUNT_TYPE_ALIAS_OBJECT:              {"ALIAS_OBJECT"},
	SAM_ACCOUNT_TYPE_NON_SECURITY_ALIAS_OBJECT: {"NON_SECURITY_ALIAS_OBJECT"},
	SAM_ACCOUNT_TYPE_USER_OBJECT:               {"USER_OBJECT", "NORMAL_USER_ACCOUNT"},
	SAM_ACCOUNT_TYPE_MACHINE_ACCOUNT:           {"MACHINE_ACCOUNT"},
	SAM_ACCOUNT_TYPE_TRUST_ACCOUNT:             {"TRUST_ACCOUNT"},
	SAM_ACCOUNT_TYPE_APP_BASIC_GROUP:           {"APP_BASIC_GROUP"},
	SAM_ACCOUNT_TYPE_APP_QUERY_GROUP:           {"APP_QUERY_GROUP"},
	SAM_ACCOUNT_TYPE_ACCOUNT_TYPE_MAX:          {"ACCOUNT_TYPE_MAX"},
	SAM_ACCOUNT_TYPE_UNKNOWN:                   {"UNKNOWN"},
}

type SamAccountType int64

func (v SamAccountType) Eval() (types []string) {
	for key, value := range samAccountTypeToString {
		if v&key != 0 {
			types = append(types, value...)
		}
	}

	return types
}

func (s SamAccountType) Strings() []string {
	if v, ok := samAccountTypeToString[s]; ok {
		return v
	}

	return samAccountTypeToString[SAM_ACCOUNT_TYPE_UNKNOWN]
}
