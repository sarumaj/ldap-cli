package objects

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

type User struct {
	AccountExpires             int64     `csv:"Expires"`
	AccountExpiryDate          time.Time `csv:"-"`
	BadPasswordDate            time.Time `csv:"-"`
	BadPasswordTime            int64     `csv:"-"`
	BadPwdCount                int       `csv:"-"`
	CommonName                 string    `ldap_attr:"cn" csv:"CN"`
	DistinguishedName          string    `csv:"DN"`
	Enabled                    bool      `csv:"Enabled"`
	GivenName                  string    `csv:"GivenName"`
	GlobalExtensionAttribute11 string    `ldap_attr:"global-ExtensionAttribute11" csv:"GlobalExtensionAttribute11"`
	GlobalExtensionAttribute22 string    `ldap_attr:"global-ExtensionAttribute22" csv:"GlobalExtensionAttribute22"`
	GlobalExtensionAttribute26 string    `ldap_attr:"global-ExtensionAttribute26" csv:"GlobalExtensionAttribute26"`
	LockedOut                  bool      `csv:"-"`
	Mail                       string    `csv:"Email"`
	MemberOf                   []string  `csv:"-"`
	MsRadiusFramedIpAddressRaw int64     `ldap_attr:"msRadiusFramedIPAddress" csv:"-"`
	MsRadiusFramedIpAddress    string    `csv:"-"`
	Name                       string    `csv:"Name"`
	ObjectCategory             string    `csv:"Category"`
	ObjectClass                []string  `csv:"-"`
	ObjectGUID                 string    `csv:"-"`
	SamAccountName             string    `csv:"sAMAccountName"`
	SamAccountType             []string  `ldap_attr:"-" csv:"-"`
	SamAccountTypeRaw          int64     `ldap_attr:"sAMAccountType" csv:"-"`
	SID                        string    `ldap_attr:"objectSid" csv:"-"`
	Surname                    string    `ldap_attr:"sn" csv:"Surname"`
	UserAccountControlRaw      int64     `ldap_attr:"userAccountControl" csv:"UserAccountControl"`
	UserAccountControl         []string  `ldap_attr:"-" csv:"-"`
	UserPrincipalName          string    `csv:"UserPrincipalName"`
}

func (u User) DN() string { return GetField(&u, "DistinguishedName") }

func (u *User) Read(raw map[string]interface{}) error {
	err := readMap(u, raw)
	if err != nil {
		return err
	}

	// place for possible implementation of custom computed properties
	if v := attributes.UserAccountControl(u.UserAccountControlRaw); v != 0 {
		u.Enabled = v&attributes.USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE == 0
		u.LockedOut = v&attributes.USER_ACCOUNT_CONTROL_LOCKOUT != 0
		u.UserAccountControl = v.Eval()
	}

	if v := u.MsRadiusFramedIpAddressRaw; v != 0 {
		u.MsRadiusFramedIpAddress = fmt.Sprintf(
			"%s.%s.%s.%s",
			strconv.FormatInt((v>>24)&0xff, 10),
			strconv.FormatInt((v>>16)&0xff, 10),
			strconv.FormatInt((v>>8)&0xff, 10),
			strconv.FormatInt((v&0xff), 10),
		)
	}

	if u.AccountExpires > 0 && u.AccountExpires < 1<<63-1 {
		u.AccountExpiryDate = util.TimeEpochBegin.Add(time.Duration(u.AccountExpires*100) * time.Nanosecond)
	}

	if u.BadPasswordTime > 0 && u.BadPasswordTime < 1<<63-1 {
		u.BadPasswordDate = util.TimeEpochBegin.Add(time.Duration(u.BadPasswordTime*100) * time.Nanosecond)
	}

	if u.ObjectGUID != "" {
		u.ObjectGUID = hexify(u.ObjectGUID)
	}

	if u.SID != "" {
		u.SID = hexify(u.SID)
	}

	u.SamAccountType = attributes.SamAccountType(u.SamAccountTypeRaw).Eval()
	return nil
}
