package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	monkey "bou.ke/monkey"
	supererrors "github.com/sarumaj/go-super/errors"
	commands "github.com/sarumaj/ldap-cli/pkg/app/commands"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestExecute(t *testing.T) {
	libutil.SkipOAT(t)

	defer monkey.Patch(os.Exit, func(code int) {
		if code > 0 {
			t.Errorf("exit code: %d", code)
		}
	}).Unpatch()

	supererrors.RegisterCallback(func(err error) {
		switch err {

		case nil, io.EOF:

		default:
			t.Error(err)

		}
	})
	apputil.Logger.SetOutput(io.Discard)

	supererrors.Except(os.Setenv("NO_COLOR", "true"))
	defer supererrors.Except(os.Unsetenv("NO_COLOR"))

	stdOut, stdErr := os.Stdout, os.Stderr
	defer func() { os.Stdout, os.Stderr = stdOut, stdErr }()

	bindParameters := []string{"--user", "cn=admin,dc=mock,dc=ad,dc=com", "--password", "admin", "--url", "ldap://localhost:389"}
	getParameters := append(bindParameters, "get", "--path", "dc=mock,dc=ad,dc=com", "--select", "*")
	editParameters := append(bindParameters, "edit", "--path", "dc=mock,dc=ad,dc=com")
	for _, tt := range []struct {
		name string
		args []string
	}{
		{"test#1", []string{"version"}},
		{"test#2", append(getParameters, "custom", "--filter", "(cn=uix00001)")},
		{"test#3", append(getParameters, "custom", "--filter", "(cn=group01)")},
		{"test#4", append(getParameters, "user", "--user-id", "uix00001")},
		{"test#5", append(getParameters, "group", "--group-id", "group01")},
		{"test#6", append(editParameters, "user", "--user-id", "uix00001", "--new-password", "new-password", "--password-attribute", "userPassword")},
		{"test#7", append(editParameters, "group", "--group-id", "group01", "--add-member", "uix00002", "--member-attribute", "memberUid")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer := supererrors.ExceptFn2(supererrors.W2(os.Pipe()))
			os.Stdout, os.Stderr = writer, writer

			t.Log("Command: ldap-cli " + strings.Join(tt.args, " "))
			commands.Execute(Version, BuildDate, tt.args...)

			supererrors.Except(writer.Close())
			buffer := bytes.NewBuffer(nil)
			_ = supererrors.ExceptFn(supererrors.W(io.Copy(buffer, reader)))
			t.Log("Output: " + buffer.String())
		})
	}
}
