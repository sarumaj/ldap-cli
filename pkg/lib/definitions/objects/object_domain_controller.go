package objects

import (
	"time"

	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type DomainController struct {
	AccountExpires        int64     `csv:"Expires"`
	AccountExpiryDate     time.Time `csv:"-"`
	Description           string    `csv:"Description"`
	DistinguishedName     string    `csv:"DN"`
	Enabled               bool      `csv:"Enabled"`
	Hostname              string    `ldap_attr:"dnsHostname" csv:"Hostname"`
	ObjectCategory        string    `csv:"Category"`
	ObjectClass           []string  `csv:"-"`
	ObjectGUID            string    `csv:"-"`
	SamAccountName        string    `csv:"sAMAccountName"`
	SamAccountType        []string  `ldap_attr:"-" csv:"-"`
	SamAccountTypeRaw     int64     `ldap_attr:"sAMAccountType" csv:"-"`
	SID                   string    `ldap_attr:"objectSid" csv:"-"`
	UserAccountControlRaw int64     `ldap_attr:"userAccountControl" csv:"UserAccountControl"`
	UserAccountControl    []string  `ldap_attr:"-" csv:"-"`
}

func (d DomainController) DN() string { return GetField(&d, "DistinguishedName") }

func (d *DomainController) Read(raw map[string]any) error {
	err := readMap(d, raw)
	if err != nil {
		return err
	}

	if d.AccountExpires > 0 && d.AccountExpires < 1<<63-1 {
		d.AccountExpiryDate = util.TimeEpochBegin.Add(time.Duration(d.AccountExpires*100) * time.Nanosecond)
	}

	// place for possible implementation of custom computed properties
	if v := attributes.UserAccountControl(d.UserAccountControlRaw); v != 0 {
		d.Enabled = v&attributes.USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE == 0
		d.UserAccountControl = v.Eval()
	}

	if d.ObjectGUID != "" {
		d.ObjectGUID = hexify(d.ObjectGUID)
	}

	if d.SID != "" {
		d.SID = hexify(d.SID)
	}

	d.SamAccountType = attributes.SamAccountType(d.SamAccountTypeRaw).Eval()
	return nil
}
