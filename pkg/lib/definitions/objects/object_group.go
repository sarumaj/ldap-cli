package objects

import "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"

type Group struct {
	CommonName        string   `ldap_attr:"cn" csv:"CN"`
	Description       string   `csv:"Description"`
	DisplayName       string   `csv:"DisplayName"`
	DistinguishedName string   `csv:"DN"`
	GroupTypeRaw      int64    `ldap_attr:"groupType" csv:"GroupType"`
	GroupType         []string `ldap_attr:"-" csv:"-"`
	MemberOf          []string `csv:"-"`
	Members           []string `ldap_attr:"member" csv:"-"`
	Name              string   `csv:"Name"`
	ObjectCategory    string   `csv:"Category"`
	ObjectClass       []string `csv:"-"`
	ObjectGUID        string   `csv:"-"`
	SamAccountName    string   `csv:"sAMAccountName"`
	SamAccountType    []string `ldap_attr:"-" csv:"-"`
	SamAccountTypeRaw int64    `ldap_attr:"sAMAccountType" csv:"-"`
	SID               string   `ldap_attr:"objectSid" csv:"-"`
}

func (g Group) DN() string { return GetField[string](&g, "DistinguishedName") }

func (g *Group) Read(raw map[string]any) error {
	err := readMap(g, raw)
	if err != nil {
		return err
	}

	// place for possible implementation of custom computed properties
	if v := attributes.GroupType(g.GroupTypeRaw); v != 0 {
		g.GroupType = v.Eval()
	}

	if g.ObjectGUID != "" {
		g.ObjectGUID = hexify(g.ObjectGUID)
	}

	if g.SID != "" {
		g.SID = hexify(g.SID)
	}

	g.SamAccountType = attributes.SamAccountType(g.SamAccountTypeRaw).Eval()

	return nil
}
