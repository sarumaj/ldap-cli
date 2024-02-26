package util

import (
	"fmt"
	"os"
	"reflect"

	keyring "github.com/99designs/keyring"
	survey "github.com/AlecAivazis/survey/v2"
)

// forked KEYCTL_PERM modes from github.com/99designs/keyring to make the config platform independent
// e.g. (KEYCTL_PERM_ALL << KEYCTL_PERM_USER) | (KEYCTL_PERM_ALL << KEYCTL_PERM_PROCESS)
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

var availableBackends = keyring.AvailableBackends()

var backendOrder = []keyring.BackendType{
	// Windows
	keyring.WinCredBackend,
	// Linux
	keyring.SecretServiceBackend,
	keyring.KWalletBackend,
	keyring.KeyCtlBackend,
	// MacOS
	keyring.KeychainBackend,
	// General
	keyring.PassBackend,
	keyring.FileBackend,
}

// Config is the configuration for the keyring
var Config = func() keyring.Config {
	cfg := &keyring.Config{
		AllowedBackends:                nil, // will be overwritten
		FileDir:                        "~/.config/ldap-cli",
		FilePasswordFunc:               passwordFunc,
		KeyCtlScope:                    "user",
		KeyCtlPerm:                     (KEYCTL_PERM_ALL << KEYCTL_PERM_USER) | (KEYCTL_PERM_ALL << KEYCTL_PERM_PROCESS),
		KeychainAccessibleWhenUnlocked: true,
		KeychainName:                   "ldap-cli",
		KeychainPasswordFunc:           nil, // will be overwritten
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

	// avoid interactive password prompt
	cfg.KeychainPasswordFunc = keyring.FixedStringPrompt("test")

	// evaluate available backends
	var backends []keyring.BackendType
	for _, backend := range backendOrder {
		// skip backend if not available
		if !func() bool {
			for _, available := range availableBackends {
				if backend == available {
					return true
				}
			}

			return false
		}() ||
			// skip file backend, since it is supposed to be always available
			backend == keyring.FileBackend {

			continue
		}

		cfg.AllowedBackends = []keyring.BackendType{backend}
		// try to open the keyring
		ring, err := keyring.Open(*cfg)
		if err != nil {
			continue
		}

		// try to set
		if err := ring.Set(keyring.Item{Key: "test", Data: []byte("test")}); err != nil {
			continue
		}

		// cleanup
		_ = ring.Remove("test")

		// backend is available
		backends = append(backends, backend)
	}

	// add file backend as last resort
	backends = append(backends, keyring.FileBackend)

	// set evaluated backends
	cfg.AllowedBackends = backends

	// reenable interactive password prompt
	cfg.KeychainPasswordFunc = passwordFunc

	return *cfg
}()

// OpenKeyring is a reference to keyring.Open (can be overwritten for testing)
var OpenKeyring = keyring.Open

// GetFromKeyring retrieves a value from the keyring
func GetFromKeyring(key string) (string, error) {
	ring, err := OpenKeyring(Config)
	if err != nil {
		return "", err
	}

	item, err := ring.Get(key)
	if err != nil {
		if ErrorIs(err, keyring.ErrKeyNotFound, os.ErrNotExist) {
			return "", nil
		}

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
		func(ans any) error {
			if reflect.Indirect(reflect.ValueOf(ans)).Len() < 12 {
				//lint:ignore ST1005 this error message should render as capitalized
				return fmt.Errorf("Password must be at least 12 characters long")
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
	ring, err := OpenKeyring(Config)
	if err != nil {
		return err
	}

	if err := ring.Remove(key); err != nil && !ErrorIs(err, keyring.ErrKeyNotFound, os.ErrNotExist) {
		return err
	}

	return nil
}

// SetToKeyring sets a value to the keyring
func SetToKeyring(key, value string) error {
	ring, err := OpenKeyring(Config)
	if err != nil {
		return err
	}

	if value == "" {
		if err := ring.Remove(key); err != nil && !ErrorIs(err, keyring.ErrKeyNotFound, os.ErrNotExist) {
			return err
		}

		return nil
	}

	return ring.Set(keyring.Item{
		Key:                         key,
		Data:                        []byte(value),
		Label:                       "LDAP CLI",
		KeychainNotTrustApplication: !Config.KeychainTrustApplication,
		KeychainNotSynchronizable:   !Config.KeychainSynchronizable,
	})
}
