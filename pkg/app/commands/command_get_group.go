package commands

import (
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var defaultGroupAttributes = attributes.Attributes{
	attributes.CommonName(),
	attributes.Description(),
	attributes.DisplayName(),
	attributes.DistinguishedName(),
	attributes.GroupType(),
	attributes.Name(),
	attributes.ObjectCategory(),
	attributes.ObjectClass(),
	attributes.SamAccountName(),
	attributes.SamAccountType(),
	attributes.UserPrincipalName(),
}

var getGroupFlags struct {
	id string
}

var getGroupCmd = func() *cobra.Command {
	getGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Get a group(s) in the directory",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"get --path \"DC=example,DC=com\" --select \"sAmAccountName,Members\" " +
			"group --group-id \"uix12345\"",
		PersistentPreRun: getGroupPersistentPreRun,
		Run:              getXRun,
	}

	flags := getGroupCmd.Flags()
	flags.StringVar(&getGroupFlags.id, "group-id", "", "Specify group ID (common name, DN or SAN)")

	return getGroupCmd
}()

func getGroupPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	getFlags.searchArguments.Attributes = append(getFlags.searchArguments.Attributes, defaultGroupAttributes...)

	var filters []filter.Filter
	if getGroupFlags.id != "" {
		filters = append(filters, filter.Or(
			filter.Filter{Attribute: attributes.SamAccountName(), Value: getGroupFlags.id},
			filter.Filter{Attribute: attributes.UserPrincipalName(), Value: getGroupFlags.id},
			filter.Filter{Attribute: attributes.Name(), Value: getGroupFlags.id},
			filter.Filter{Attribute: attributes.DistinguishedName(), Value: getGroupFlags.id},
		))
	}

	getFlags.searchArguments.Filter = filter.And(filter.IsGroup(), filters...)
}
