package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var defaultCustomAttributes = attributes.Attributes{
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
	filterString string
}

var getCustomCmd = func() *cobra.Command {
	getCustomCmd := &cobra.Command{
		Use:              "custom",
		Short:            "Get an arbitrary directory object",
		Example:          "ldap-cli get custom",
		PersistentPreRun: getCustomPersistentPreRun,
		Run:              getXRun,
	}

	flags := getCustomCmd.Flags()
	flags.StringVar(&getCustomFlags.filterString, "filter", "", "Provide custom LDAP query filter")

	return getCustomCmd
}()

func getCustomPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	getFlags.searchArguments.Attributes = append(getFlags.searchArguments.Attributes, defaultCustomAttributes...)

	getFlags.searchArguments.Filter = *supererrors.ExceptFn(supererrors.W(filter.ParseRaw(getCustomFlags.filterString)))
}
