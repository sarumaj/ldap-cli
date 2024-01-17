package util

import (
	"os/exec"
	"testing"
)

func TestKeyringFlow(t *testing.T) {
	if _, err := exec.LookPath("pass"); err != nil {
		t.Skipf(`Skipping keyring test, since it requires "pass" to be installed`)
	}

	if err := SetToKeyring("test", "12345"); err != nil {
		t.Errorf(`SetToKeyring("test", "12345") failed: %v`, err)
	}

	got, err := GetFromKeyring("test")
	if err != nil {
		t.Errorf(`GetFromKeyring("test") failed: %v`, err)
		return
	}

	if got != "12345" {
		t.Errorf(`GetFromKeyring("test") failed: got: %q, want: %q`, got, "12345")
	}

	if err := RemoveFromKeyRing("test"); err != nil {
		t.Errorf(`RemoveFromKeyRing("test") failed: %v`, err)
	}
}
