package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	logrus "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var loggerEntry = apputil.Logger.WithFields(logrus.Fields{"mod": "commands"})

var rootCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ldap-cli",
		Short: "ldap-cli is cross-platform compatible client application based on the lightweight directory access control (LDAP)",
		Run: func(cmd *cobra.Command, _ []string) {
			supererrors.Except(cmd.Help())
		},
		Example: "ldap-cli <subcommand>",
		Version: internalVersion,
	}

	cmd.AddCommand(versionCmd)

	return cmd
}()

// Execute executes the root command.
func Execute(version, buildDate string) {
	internalVersion, internalBuildDate = version, buildDate

	loggerEntry.Debugf("Version: %s, build date: %s, executable path: %s", internalVersion, internalBuildDate, apputil.GetExecutablePath())

	if err := rootCmd.Execute(); err != nil {
		loggerEntry.Debugf("Execution failed: %v", err)
	}
}
