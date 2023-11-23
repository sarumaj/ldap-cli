package main

import (
	"os"
	"runtime/pprof"

	supererrors "github.com/sarumaj/go-super/errors"
	commands "github.com/sarumaj/ldap-cli/pkg/app/commands"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
)

// Version holds the application version.
// It gets filled automatically at build time.
var Version = "v0.0.0"

// BuildDate holds the date and time at which the application was build.
// It gets filled automatically at build time.
var BuildDate = "0000-00-00 00:00:00 UTC"

func main() {
	apputil.Logger.Debugf("version: %q, build date: %q", Version, BuildDate)

	f := supererrors.ExceptFn(supererrors.W(os.OpenFile("profile.prof", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)))
	defer f.Close()

	supererrors.Except(pprof.StartCPUProfile(f))
	defer pprof.StopCPUProfile()

	commands.Execute(Version, BuildDate)
}
