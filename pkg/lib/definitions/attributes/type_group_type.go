package attributes

import "slices"

// https://docs.microsoft.com/en-us/windows/win32/adschema/a-grouptype
const (
	GROUP_TYPE_CREATED_BY_SYSTEM FlagsetGroupType = 0x00000001 // Specifies a group that is created by the system.
	GROUP_TYPE_GLOBAL            FlagsetGroupType = 0x00000002 // Specifies a group with global scope.
	GROUP_TYPE_LOCAL             FlagsetGroupType = 0x00000004 // Specifies a group with domain local scope.
	GROUP_TYPE_UNIVERSAL         FlagsetGroupType = 0x00000008 // Specifies a group with universal scope.
	GROUP_TYPE_APP_BASIC         FlagsetGroupType = 0x00000010 // Specifies an APP_BASIC group for Windows Server Authorization Manager.
	GROUP_TYPE_APP_QUERY         FlagsetGroupType = 0x00000020 // Specifies an APP_QUERY group for Windows Server Authorization Manager.
	GROUP_TYPE_SECURITY          FlagsetGroupType = 0x80000000 // Specifies a security group. If this flag is not set, then the group is a distribution group.
	GROUP_TYPE_DISTRIBUTION      FlagsetGroupType = ^GROUP_TYPE_SECURITY
)

var groupTypeToString = map[FlagsetGroupType]string{
	GROUP_TYPE_CREATED_BY_SYSTEM: "CREATED_BY_SYSTEM",
	GROUP_TYPE_GLOBAL:            "GLOBAL",
	GROUP_TYPE_LOCAL:             "LOCAL",
	GROUP_TYPE_UNIVERSAL:         "UNIVERSAL",
	GROUP_TYPE_APP_BASIC:         "APP_BASIC",
	GROUP_TYPE_APP_QUERY:         "APP_QUERY",
	GROUP_TYPE_SECURITY:          "SECURITY",
	GROUP_TYPE_DISTRIBUTION:      "DISTRIBUTION",
}

type FlagsetGroupType uint32

func (v FlagsetGroupType) Eval() (types []string) {
	for key, value := range groupTypeToString {
		if v&key == key {
			types = append(types, value)
		}
	}

	if v&GROUP_TYPE_SECURITY != GROUP_TYPE_SECURITY && len(types) > 0 {
		types = append(types, GROUP_TYPE_DISTRIBUTION.String())
	}

	slices.Sort(types)
	return types
}

func (g FlagsetGroupType) String() string {
	if v, ok := groupTypeToString[g]; ok {
		return v
	}

	return "UNKNOWN"
}
