package commands

import (
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

// Default attributes for search query
var defaultGroupGetAttributes = attributes.Attributes{
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

// Command options
var getGroupFlags struct {
	id string `flag:"group-id"`
}

// "group" command
var getGroupCmd = func() *cobra.Command {
	getGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Get a group(s) in the directory",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"get --path \"DC=example,DC=com\" --select \"sAmAccountName,Members\" " +
			"group --group-id \"uix12345\"",
		PersistentPreRun: getGroupPersistentPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			getChildCommandRun(cmd, args)
		},
	}

	flags := getGroupCmd.Flags()
	flags.StringVar(&getGroupFlags.id, "group-id", "", "Specify group ID (common name, DN or SAN)")

	return getGroupCmd
}()

// Runs always prior to "run"
func getGroupPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getGroupPersistentPreRun"})
	logger.Trace("Executing")

	apputil.AskID(cmd, "group-id", &getGroupFlags.id, &getFlags.searchArguments)

	if len(getFlags.searchArguments.Attributes) == 0 {
		getFlags.searchArguments.Attributes.Append(defaultGroupGetAttributes...)
	}
	logger.WithField("searchArguments.Attributes", getFlags.searchArguments.Attributes).Trace("Set")

	var filters []filter.Filter
	if getGroupFlags.id != "" {
		filters = append(filters, filter.ByID(getGroupFlags.id))
	}

	getFlags.searchArguments.Filter = filter.And(filter.IsGroup(), filters...)
	logger.WithField("searchArguments.Filter", getFlags.searchArguments.Filter.String()).Trace("Set")
}
