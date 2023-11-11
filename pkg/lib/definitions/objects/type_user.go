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
	BadPasswordCount           int       `ldap_attr:"badPwdCount" csv:"-"`
	BadPasswordTime            time.Time `ldap_attr:"-"`
	BadPasswordTimeRaw         int64     `ldap_attr:"badPasswordTime" csv:"-"`
	CountryCode                string    `ldap_attr:"c"`
	CountryName                string    `ldap_attr:"co"`
	CommonName                 string    `ldap_attr:"cn" csv:"CN"`
	Company                    string
	DistinguishedName          string    `csv:"DN"`
	Division                   string    `ldap_attr:"division"`
	Enabled                    bool      `csv:"Enabled"`
	GivenName                  string    `csv:"GivenName"`
	GlobalExtensionAttribute11 string    `ldap_attr:"global-extensionAttribute11" csv:"GlobalExtensionAttribute11"`
	GlobalExtensionAttribute22 string    `ldap_attr:"global-extensionAttribute22" csv:"GlobalExtensionAttribute22"`
	GlobalExtensionAttribute26 string    `ldap_attr:"global-extensionAttribute26" csv:"GlobalExtensionAttribute26"`
	Location                   string    `ldap_attr:"l"`
	LockedOut                  bool      `csv:"-"`
	Mail                       string    `csv:"Email"`
	MemberOf                   []string  `csv:"-"`
	MsRadiusFramedIpAddressRaw int64     `ldap_attr:"msRadiusFramedIPAddress" csv:"-"`
	MsRadiusFramedIpAddress    string    `csv:"-"`
	Name                       string    `csv:"Name"`
	ObjectCategory             string    `csv:"Category"`
	ObjectClass                []string  `csv:"-"`
	ObjectGUID                 string    `csv:"-"`
	PasswordLastSet            time.Time `ldap_attr:"-"`
	PasswordLastSetRaw         int64     `ldap_attr:"pwdLastSet" csv:"-"`
	SamAccountName             string    `csv:"sAMAccountName"`
	SamAccountType             []string  `ldap_attr:"-" csv:"-"`
	SamAccountTypeRaw          int64     `ldap_attr:"sAMAccountType" csv:"-"`
	SID                        string    `ldap_attr:"objectSid" csv:"-"`
	Surname                    string    `ldap_attr:"sn" csv:"Surname"`
	UserAccountControl         []string  `ldap_attr:"-" csv:"-"`
	UserAccountControlRaw      int64     `ldap_attr:"userAccountControl" csv:"UserAccountControl"`
	UserPrincipalName          string    `csv:"UserPrincipalName"`
}

func (u User) DN() string { return GetField[string](&u, "DistinguishedName") }

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
		u.AccountExpiryDate = util.TimeAfter1601(u.AccountExpires)
	}

	if u.BadPasswordTimeRaw > 0 && u.BadPasswordTimeRaw < 1<<63-1 {
		u.BadPasswordTime = util.TimeAfter1601(u.AccountExpires)
	}

	if u.PasswordLastSetRaw > 0 && u.PasswordLastSetRaw < 1<<63-1 {
		u.PasswordLastSet = util.TimeAfter1601(u.PasswordLastSetRaw)
	}

	u.SamAccountType = attributes.SAMAccountType(u.SamAccountTypeRaw).Eval()
	return nil
}
