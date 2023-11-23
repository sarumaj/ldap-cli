package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var editCustomFlags = &struct {
	filterString string
}{}

var editCustomCmd = func() *cobra.Command {
	editCustomCmd := &cobra.Command{
		Use:              "custom",
		Short:            "Get an arbitrary directory object to edit",
		Example:          "ldap-cli edit custom",
		PersistentPreRun: editCustomPersistentPreRun,
		PreRun:           editChildCommandPreRun,
		Run:              editCustomRun,
		PostRun:          editChildCommandPostRun,
	}

	flags := editCustomCmd.Flags()
	flags.StringVar(&editCustomFlags.filterString, "filter", "", "Provide custom LDAP query filter")

	return editCustomCmd
}()

func editCustomPersistentPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editCustomPersistentPreRun"})
	logger.Debug("Executing")

	editFlags.searchArguments.Filter = *supererrors.ExceptFn(supererrors.W(filter.ParseRaw(editCustomFlags.filterString)))
	logger.WithField("searchArguments.Filter", editFlags.searchArguments.Filter.String()).Debug("Set")
}

func editCustomRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editCustomRun"})
	logger.Debug("Executing")

	requests := editFlags.requests
	_ = supererrors.ExceptFn(supererrors.W(apputil.AskLDAPDataInterchangeFormat(requests, editFlags.editor)))
	editFlags.requests = requests
	logger.WithField("editor", editFlags.editor).Debug("Asked")
}
