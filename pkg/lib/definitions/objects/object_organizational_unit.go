package objects

type OrganizationalUnit struct {
	CountryCode        string `ldap_attr:"c" csv:"CountryCode"`
	CountryName        string `ldap_attr:"co" csv:"CountryName"`
	Description        string
	DistinguishedName  string `csv:"DN"`
	Name               string
	ObjectCategory     string   `csv:"Category"`
	ObjectClass        []string `csv:"-"`
	ObjectGUID         string   `csv:"-"`
	OrganizationalUnit []string `ldap_attr:"ou" csv:"OrganizationalUnit"`
}

func (o OrganizationalUnit) DN() string { return GetField(&o, "DistinguishedName") }

func (o *OrganizationalUnit) Read(raw map[string]any) error {
	err := readMap(o, raw)
	if err != nil {
		return err
	}

	// place for possible implementation of custom computed properties
	if o.ObjectGUID != "" {
		o.ObjectGUID = hexify(o.ObjectGUID)
	}

	return nil
}
