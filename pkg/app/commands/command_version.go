package commands

import (
	"fmt"

	semver "github.com/blang/semver"
	color "github.com/fatih/color"
	selfupdate "github.com/rhysd/go-github-selfupdate/selfupdate"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
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

// Command options
var versionFlags struct {
	update bool
}

// "version" command
var versionCmd = func() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:     "version",
		Short:   "Display version information",
		Example: "ldap-cli version",
		Run:     versionRun,
	}

	flags := versionCmd.Flags()
	flags.BoolVar(&versionFlags.update, "update", false, "Update application to the newest version")

	return versionCmd
}()

// Check app version or/and update to the latest
func versionRun(*cobra.Command, []string) {
	current := supererrors.ExceptFn(supererrors.W(semver.ParseTolerant(internalVersion)))
	latest, found, err := selfupdate.DetectLatest(remoteRepository)

	var vSuffix string
	switch {
	case err == nil && (!found || latest.Version.LTE(current)):
		vSuffix = " (latest)"

	case err == nil && (found || latest.Version.GT(current)):
		if versionFlags.update {
			up := supererrors.ExceptFn(supererrors.W(selfupdate.NewUpdater(selfupdate.Config{
				Validator: &selfupdate.SHA2Validator{},
			})))
			_ = supererrors.ExceptFn(supererrors.W(up.UpdateSelf(current, remoteRepository)))
			_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(
				apputil.Stdout(),
				apputil.PrintColors(color.HiGreenString, "Successfully updated from %s to %s", current, latest.Version),
			)))
			return

		} else {
			vSuffix = " (newer version available: " + latest.Version.String() + ", run \"ldap-cli version --update\" to update)"

		}

	}

	fmt.Println("test")
	_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), apputil.PrintColors(color.CyanString, "Version: %s", internalVersion+vSuffix))))
	_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), apputil.PrintColors(color.CyanString, "Built at: %s", internalBuildDate))))
	_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), apputil.PrintColors(color.CyanString, "Executable path: %s", libutil.GetExecutablePath()))))
}
