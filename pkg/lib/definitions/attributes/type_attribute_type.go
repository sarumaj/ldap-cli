package attributes

const (
	TypeBool               AttributeType = "Bool"
	TypeDecimal            AttributeType = "Decimal"
	TypeGroupType          AttributeType = "GroupType"
	TypeHexString          AttributeType = "HexString"
	TypeInt                AttributeType = "Int"
	TypeIPv4Address        AttributeType = "IPv4Address"
	TypeRaw                AttributeType = "TypeRaw"
	TypeSAMaccountType     AttributeType = "SAMaccountType"
	TypeString             AttributeType = "String"
	TypeStringSlice        AttributeType = "StringSlice"
	TypeTime               AttributeType = "Time"
	TypeUserAccountControl AttributeType = "UserAccountControl"
)

type AttributeType string
