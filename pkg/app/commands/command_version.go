package commands

import (
	"fmt"

	semver "github.com/blang/semver"
	color "github.com/fatih/color"
	selfupdate "github.com/rhysd/go-github-selfupdate/selfupdate"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	cobra "github.com/spf13/cobra"
)

// Address of remote repository where the newest version of gh-gr is released.
const remoteRepository = "sarumaj/ldap-cli"

// Version holds the application version.
// It gets filled automatically at build time.
var internalVersion string

// BuildDate holds the date and time at which the application was build.
// It gets filled automatically at build time.
var internalBuildDate string

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Display version information",
	Example: "ldap-cli version",
	Run: func(*cobra.Command, []string) {
		current := supererrors.ExceptFn(supererrors.W(semver.ParseTolerant(internalVersion)))
		latest, found, err := selfupdate.DetectLatest(remoteRepository)

		var vSuffix string
		switch {
		case err == nil && (!found || latest.Version.LTE(current)):
			vSuffix = " (latest)"

		case err == nil && (found || latest.Version.GT(current)):
			vSuffix = " (newer version available: " + latest.Version.String() + ", run \"gh extension upgrade gr\" to update)"

		}

		_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), apputil.CheckColors(color.CyanString, "Version: %s", internalVersion+vSuffix))))
		_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), apputil.CheckColors(color.CyanString, "Built at: %s", internalBuildDate))))
		_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), apputil.CheckColors(color.CyanString, "Executable path: %s", apputil.GetExecutablePath()))))
	},
}
