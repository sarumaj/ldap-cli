package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var defaultCustomGetAttributes = attributes.Attributes{
	attributes.CommonName(),
	attributes.DisplayName(),
	attributes.DistinguishedName(),
	attributes.Name(),
	attributes.ObjectCategory(),
	attributes.ObjectClass(),
	attributes.SamAccountName(),
	attributes.SamAccountType(),
}

var getCustomFlags struct {
	filterString string `flag:"filter"`
}

var getCustomCmd = func() *cobra.Command {
	getCustomCmd := &cobra.Command{
		Use:              "custom",
		Short:            "Get an arbitrary directory object",
		Example:          "ldap-cli get custom",
		PersistentPreRun: getCustomPersistentPreRun,
		Run:              getChildCommandRun,
	}

	flags := getCustomCmd.Flags()
	flags.StringVar(&getCustomFlags.filterString, "filter", "", "Provide custom LDAP query filter")

	return getCustomCmd
}()

func getCustomPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getCustomPersistentPreRun"})
	logger.Debug("Executing")

	if getCustomFlags.filterString == "" {
		var args []string
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(cmd, "filter", &args, false)))
		supererrors.Except(cmd.ParseFlags(args))
		getFlags.searchArguments.Filter = *supererrors.ExceptFn(supererrors.W(filter.ParseRaw(getCustomFlags.filterString)))
		logger.WithField("searchArguments.Filter", getFlags.searchArguments.Filter).Debug("Asked")
	}

	if len(getFlags.searchArguments.Attributes) == 0 {
		getFlags.searchArguments.Attributes.Append(defaultCustomGetAttributes...)
	}
	logger.WithField("searchArguments.Attributes", getFlags.searchArguments.Attributes).Debug("Set")

	getFlags.searchArguments.Filter = *supererrors.ExceptFn(supererrors.W(filter.ParseRaw(getCustomFlags.filterString)))
	logger.WithField("searchArguments.Filter", getFlags.searchArguments.Filter.String()).Debug("Set")
}
