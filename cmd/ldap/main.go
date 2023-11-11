package main

import (
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
	commands.Execute(Version, BuildDate)
}
