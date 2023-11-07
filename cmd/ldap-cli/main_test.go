package main

import (
	"io"
	"testing"

	supererrors "github.com/sarumaj/go-super/errors"
	"github.com/sarumaj/ldap-cli/pkg/app/commands"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func TestExecute(t *testing.T) {
	libutil.SkipOAT(t)

	apputil.Logger.SetOutput(io.Discard)
	supererrors.RegisterCallback(func(err error) { t.Error(err) })

	for _, tt := range []struct {
		name string
		args []string
	}{
		// TODO
	} {
		t.Run(tt.name, func(t *testing.T) {
			commands.Execute(Version, BuildDate, tt.args...)
		})
	}
}
