package main

import (
	"testing"

	supererrors "github.com/sarumaj/go-super/errors"
	commands "github.com/sarumaj/ldap-cli/pkg/app/commands"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestExecute(t *testing.T) {
	libutil.SkipOAT(t)

	//apputil.Logger.SetOutput(io.Discard)
	supererrors.RegisterCallback(func(err error) { t.Error(err) })

	for _, tt := range []struct {
		name string
		args []string
	}{
		{"test#1", []string{"version"}},
		{"test#2", []string{"--debug", "--user", "cn=admin,dc=mock,dc=ad,dc=com", "--password", "admin", "--url", "ldap://localhost:389", "get", "--path", "dc=mock,dc=ad,dc=com", "user", "--user-id", "cn=admin,dc=mock,dc=ad,dc=com"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commands.Execute(Version, BuildDate, tt.args...)
		})
	}
}
