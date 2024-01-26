package attributes

import "slices"

// userAccountControl implemented corresponding to: [https://docs.microsoft.com/en-us/windows/win32/adschema/a-useraccountcontrol].
const (
	USER_ACCOUNT_CONTROL_LOGON_SCRIPT                           FlagsetUserAccountControl = 0x00000001 // The logon script is executed.
	USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE                        FlagsetUserAccountControl = 0x00000002 // The user account is disabled.
	USER_ACCOUNT_CONTROL_HOMEDIR_REQUIRED                       FlagsetUserAccountControl = 0x00000008 // The home directory is required.
	USER_ACCOUNT_CONTROL_LOCKOUT                                FlagsetUserAccountControl = 0x00000010 // The account is currently locked out.
	USER_ACCOUNT_CONTROL_PASSWORD_NOTREQD                       FlagsetUserAccountControl = 0x00000020 // No password is required.
	USER_ACCOUNT_CONTROL_PASSWORD_CANT_CHANGE                   FlagsetUserAccountControl = 0x00000040 // The user cannot change the password.
	USER_ACCOUNT_CONTROL_ENCRYPTED_TEXT_PASSWORD_ALLOWED        FlagsetUserAccountControl = 0x00000080 // The user can send an encrypted password.
	USER_ACCOUNT_CONTROL_TEMP_DUPLICATE_ACCOUNT                 FlagsetUserAccountControl = 0x00000100 // This is an account for users whose primary account is in another domain. This account provides user access to this domain, but not to any domain that trusts this domain. Also known as a local user account.
	USER_ACCOUNT_CONTROL_NORMAL_ACCOUNT                         FlagsetUserAccountControl = 0x00000200 // This is a default account type that represents a typical user.
	USER_ACCOUNT_CONTROL_INTERDOMAIN_TRUST_ACCOUNT              FlagsetUserAccountControl = 0x00000800 // This is a permit to trust account for a system domain that trusts other domains.
	USER_ACCOUNT_CONTROL_WORKSTATION_TRUST_ACCOUNT              FlagsetUserAccountControl = 0x00001000 // This is a computer account for a computer that is a member of this domain.
	USER_ACCOUNT_CONTROL_SERVER_TRUST_ACCOUNT                   FlagsetUserAccountControl = 0x00002000 // This is a computer account for a system backup domain controller that is a member of this domain.
	USER_ACCOUNT_CONTROL_DONT_EXPIRE_PASSWD                     FlagsetUserAccountControl = 0x00010000 // The password for this account will never expire.
	USER_ACCOUNT_CONTROL_MNS_LOGON_ACCOUNT                      FlagsetUserAccountControl = 0x00020000 // This is an MNS logon account.
	USER_ACCOUNT_CONTROL_SMARTCARD_REQUIRED                     FlagsetUserAccountControl = 0x00040000 // The user must log on using a smart card.
	USER_ACCOUNT_CONTROL_TRUSTED_FOR_DELEGATION                 FlagsetUserAccountControl = 0x00080000 // The service account (user or computer account), under which a service runs, is trusted for Kerberos delegation. Any such service can impersonate a client requesting the service.
	USER_ACCOUNT_CONTROL_NOT_DELEGATED                          FlagsetUserAccountControl = 0x00100000 // The security context of the user will not be delegated to a service even if the service account is set as trusted for Kerberos delegation.
	USER_ACCOUNT_CONTROL_USE_DES_KEY_ONLY                       FlagsetUserAccountControl = 0x00200000 // Restrict this principal to use only Data Encryption Standard (DES) encryption types for keys.
	USER_ACCOUNT_CONTROL_DONT_REQUIRE_PREAUTH                   FlagsetUserAccountControl = 0x00400000 // This account does not require Kerberos pre-authentication for logon.
	USER_ACCOUNT_CONTROL_PASSWORD_EXPIRED                       FlagsetUserAccountControl = 0x00800000 // The user password has expired. This flag is created by the system using data from the Pwd-Last-Set attribute and the domain policy.
	USER_ACCOUNT_CONTROL_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION FlagsetUserAccountControl = 0x01000000 // The account is enabled for delegation. This is a security-sensitive setting; accounts with this option enabled should be strictly controlled. This setting enables a service running under the account to assume a client identity and authenticate as that user to other remote servers on the network.
)

// userAccountControlToString is a map of user account control flags to their string representation (flags are multi-valued).
var userAccountControlToString = map[FlagsetUserAccountControl]string{
	USER_ACCOUNT_CONTROL_LOGON_SCRIPT:                           "LOGON_SCRIPT",
	USER_ACCOUNT_CONTROL_ACCOUNT_DISABLE:                        "ACCOUNT_DISABLE",
	USER_ACCOUNT_CONTROL_HOMEDIR_REQUIRED:                       "HOMEDIR_REQUIRED",
	USER_ACCOUNT_CONTROL_LOCKOUT:                                "LOCKOUT",
	USER_ACCOUNT_CONTROL_PASSWORD_NOTREQD:                       "PASSWORD_NOTREQD",
	USER_ACCOUNT_CONTROL_PASSWORD_CANT_CHANGE:                   "PASSWORD_CANT_CHANGE",
	USER_ACCOUNT_CONTROL_ENCRYPTED_TEXT_PASSWORD_ALLOWED:        "ENCRYPTED_TEXT_PASSWORD_ALLOWED",
	USER_ACCOUNT_CONTROL_TEMP_DUPLICATE_ACCOUNT:                 "TEMP_DUPLICATE_ACCOUNT",
	USER_ACCOUNT_CONTROL_NORMAL_ACCOUNT:                         "NORMAL_ACCOUNT",
	USER_ACCOUNT_CONTROL_INTERDOMAIN_TRUST_ACCOUNT:              "INTERDOMAIN_TRUST_ACCOUNT",
	USER_ACCOUNT_CONTROL_WORKSTATION_TRUST_ACCOUNT:              "WORKSTATION_TRUST_ACCOUNT",
	USER_ACCOUNT_CONTROL_SERVER_TRUST_ACCOUNT:                   "SERVER_TRUST_ACCOUNT",
	USER_ACCOUNT_CONTROL_DONT_EXPIRE_PASSWD:                     "DONT_EXPIRE_PASSWD",
	USER_ACCOUNT_CONTROL_MNS_LOGON_ACCOUNT:                      "MNS_LOGON_ACCOUNT",
	USER_ACCOUNT_CONTROL_SMARTCARD_REQUIRED:                     "SMARTCARD_REQUIRED",
	USER_ACCOUNT_CONTROL_TRUSTED_FOR_DELEGATION:                 "TRUSTED_FOR_DELEGATION",
	USER_ACCOUNT_CONTROL_NOT_DELEGATED:                          "NOT_DELEGATED",
	USER_ACCOUNT_CONTROL_USE_DES_KEY_ONLY:                       "USE_DES_KEY_ONLY",
	USER_ACCOUNT_CONTROL_DONT_REQUIRE_PREAUTH:                   "DONT_REQUIRE_PREAUTH",
	USER_ACCOUNT_CONTROL_PASSWORD_EXPIRED:                       "PASSWORD_EXPIRED",
	USER_ACCOUNT_CONTROL_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION: "TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION",
}

// FlagsetUserAccountControl is a set of flags for the userAccountControl attribute (multi-valued).
type FlagsetUserAccountControl uint32

// Eval returns a list of strings representing the flags set in the userAccountControl attribute.
func (v FlagsetUserAccountControl) Eval() (controls []string) {
	for key, value := range userAccountControlToString {
		if v&key == key {
			controls = append(controls, value)
		}
	}

	slices.Sort(controls)
	return controls
}
