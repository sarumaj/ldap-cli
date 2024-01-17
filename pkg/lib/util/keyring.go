package util

import (
	"fmt"

	keyring "github.com/99designs/keyring"
	survey "github.com/AlecAivazis/survey/v2"
)

// forked KEYCTL_PERM modes from github.com/99designs/keyring to make the config platform independent
const (
	KEYCTL_PERM_VIEW = uint32(1 << iota)
	KEYCTL_PERM_READ
	KEYCTL_PERM_WRITE
	KEYCTL_PERM_SEARCH
	KEYCTL_PERM_LINK
	KEYCTL_PERM_SETATTR
	KEYCTL_PERM_ALL = uint32((1 << iota) - 1)
)

const (
	KEYCTL_PERM_OTHERS uint32 = iota * 8
	KEYCTL_PERM_GROUP
	KEYCTL_PERM_USER
	KEYCTL_PERM_PROCESS
)

// Config is the configuration for the keyring
var Config = keyring.Config{
	AllowedBackends:                keyring.AvailableBackends(),
	FileDir:                        "~/.config/ldap-cli",
	FilePasswordFunc:               passwordFunc,
	KeyCtlScope:                    "user",
	KeyCtlPerm:                     (KEYCTL_PERM_ALL << KEYCTL_PERM_USER) | (KEYCTL_PERM_ALL << KEYCTL_PERM_PROCESS),
	KeychainAccessibleWhenUnlocked: true,
	KeychainName:                   "ldap-cli",
	KeychainPasswordFunc:           passwordFunc,
	KeychainSynchronizable:         false,
	KeychainTrustApplication:       true,
	KWalletAppID:                   "ldap-cli",
	KWalletFolder:                  "ldap-cli",
	LibSecretCollectionName:        "ldap-cli",
	PassCmd:                        "pass",
	PassDir:                        "~/.password-store",
	PassPrefix:                     "ldap-cli.",
	ServiceName:                    "ldap-cli",
	WinCredPrefix:                  "ldap-cli.",
}

// GetFromKeyring retrieves a value from the keyring
func GetFromKeyring(key string) (string, error) {
	ring, err := keyring.Open(Config)
	if err != nil {
		return "", err
	}

	item, err := ring.Get(key)
	if err != nil {
		return "", err
	}

	return string(item.Data), nil
}

// passwordFunc is a helper function to ask for a password
func passwordFunc(s string) (string, error) {
	var got string
	prompt := &survey.Password{Message: "Please, provide a password to secure your credentials"}
	err := survey.AskOne(prompt, &got, survey.WithValidator(survey.ComposeValidators(
		survey.Required,
		func(ans interface{}) error {
			if str, ok := ans.(string); !ok || len(str) < 12 {
				return fmt.Errorf("password must be at least 12 characters long")
			}

			return nil
		},
	)))

	if err != nil {
		return "", err
	}

	return got, nil
}

// RemoveFromKeyRing removes a value from the keyring
func RemoveFromKeyRing(key string) error {
	ring, err := keyring.Open(Config)
	if err != nil {
		return err
	}

	return ring.Remove(key)
}

// SetToKeyring sets a value to the keyring
func SetToKeyring(key, value string) error {
	ring, err := keyring.Open(Config)
	if err != nil {
		return err
	}

	return ring.Set(keyring.Item{Key: key, Data: []byte(value), Label: "LDAP CLI"})
}
