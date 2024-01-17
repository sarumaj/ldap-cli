package util

import (
	"fmt"

	keyring "github.com/99designs/keyring"
	survey "github.com/AlecAivazis/survey/v2"
)

var Config = keyring.Config{
	AllowedBackends:                keyring.AvailableBackends(),
	FileDir:                        "~/.config/ldap-cli",
	FilePasswordFunc:               passwordFunc,
	KeyCtlScope:                    "user",
	KeyCtlPerm:                     keyring.KEYCTL_PERM_USER,
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

func passwordFunc(s string) (string, error) {
	var got string
	err := survey.AskOne(&survey.Password{Message: s}, &got, survey.WithValidator(survey.ComposeValidators(survey.Required, func(ans interface{}) error {
		if str, ok := ans.(string); !ok || len(str) < 12 {
			return fmt.Errorf("password must be at least 12 characters long")
		}

		return nil
	})))

	if err != nil {
		return "", err
	}

	return got, nil
}

func RemoveFromKeyRing(key string) error {
	ring, err := keyring.Open(Config)
	if err != nil {
		return err
	}

	return ring.Remove(key)
}

func SetToKeyring(key, value string) error {
	ring, err := keyring.Open(Config)
	if err != nil {
		return err
	}

	return ring.Set(keyring.Item{Key: key, Data: []byte(value), Label: "LDAP CLI"})
}
