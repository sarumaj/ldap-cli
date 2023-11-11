package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var searchCustomCmd = func() *cobra.Command {
	searchCustomCmd := &cobra.Command{
		Use:     "custom",
		Short:   "Search an arbitrary directory object",
		Example: "ldap-cli search custom",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			parent := cmd.Parent()
			parent.PersistentPreRun(parent, nil)

			searchArguments.Attributes = append(
				searchArguments.Attributes,
				attributes.CommonName(),
				attributes.DisplayName(),
				attributes.DistinguishedName(),
				attributes.Name(),
				attributes.ObjectCategory(),
				attributes.ObjectClass(),
				attributes.SamAccountName(),
				attributes.SamAccountType(),
			)

			searchArguments.Filter = *supererrors.ExceptFn(supererrors.W(filter.ParseRaw(filterString)))
		},
		Run: performSearch,
	}

	flags := searchCustomCmd.Flags()
	flags.StringVar(&filterString, "filter", "", "Provide custom LDAP query filter")

	return searchCustomCmd
}()
