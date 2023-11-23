package attributes

const (
	TypeBool               Type = "Bool"
	TypeDecimal            Type = "Decimal"
	TypeGroupType          Type = "GroupType"
	TypeHexString          Type = "HexString"
	TypeInt                Type = "Int"
	TypeIPv4Address        Type = "IPv4Address"
	TypeRaw                Type = "TypeRaw"
	TypeSAMaccountType     Type = "SAMaccountType"
	TypeString             Type = "String"
	TypeStringSlice        Type = "StringSlice"
	TypeTime               Type = "Time"
	TypeUserAccountControl Type = "UserAccountControl"
)

type Type string
