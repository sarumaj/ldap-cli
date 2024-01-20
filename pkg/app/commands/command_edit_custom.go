package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	cobra "github.com/spf13/cobra"
)

// Command options
var editCustomFlags struct {
	filterString string `flag:"filter"`
}

// "custom" command
var editCustomCmd = func() *cobra.Command {
	editCustomCmd := &cobra.Command{
		Use:   "custom",
		Short: "Edit an arbitrary directory object",
		Long:  "Edit an arbitrary directory object.\n\n",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" edit " +
			"custom --filter \"(cn=commonName)\"",
		PersistentPreRun: editCustomPersistentPreRun,
		PreRun:           editChildCommandPreRun,
		Run:              editCustomRun,
		PostRun:          editChildCommandPostRun,
	}

	apputil.HelpAliases(&editCustomCmd.Long)

	flags := editCustomCmd.Flags()
	flags.StringVar(&editCustomFlags.filterString, "filter", "", "Provide custom LDAP query filter")

	return editCustomCmd
}()

// Runs always prior to "run"
func editCustomPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editCustomPersistentPreRun"})
	logger.Trace("Executing")

	apputil.AskFilterString(cmd, "filter", &editCustomFlags.filterString, &editFlags.searchArguments)
	logger.WithField("searchArguments.Filter", editFlags.searchArguments.Filter.String()).Trace("Set")
}

// Actual "run" prepares a modify request
func editCustomRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editCustomRun"})
	logger.Trace("Executing")

	requests := editFlags.requests
	_ = supererrors.ExceptFn(supererrors.W(apputil.AskLDAPDataInterchangeFormat(requests, editFlags.editor)))
	editFlags.requests = requests
	logger.WithField("editor", editFlags.editor).Trace("Asked")
}
