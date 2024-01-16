package util

import (
	keyring "github.com/99designs/keyring"
)

const Service = "ldap-cli"

var Config = keyring.Config{
	ServiceName:     Service,
	AllowedBackends: keyring.AvailableBackends(),
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
