package commands

import (
	"fmt"

	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

// Default attributes for search query
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

// Command options
var getCustomFlags struct {
	filterString string `flag:"filter"`
}

// "custom" command
var getCustomCmd = func() *cobra.Command {
	getCustomCmd := &cobra.Command{
		Use:   "custom",
		Short: "Get arbitrary directory object(s)",
		Long: "Get  arbitrary directory object(s).\n\n" +
			"Filter option supports following interpolations:\n",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"get --path \"DC=example,DC=com\" --select \"sAmAccountName,AccountExpires\" " +
			"custom --filter \"(&(cn=commonName)(memberof:${RECURSIVE}:=groupName))\"",
		PersistentPreRun: getCustomPersistentPreRun,
		Run:              getChildCommandRun,
	}

	for _, alias := range filter.ListAliases() {
		getCustomCmd.Long += fmt.Sprintf(" - %-24s: %s\n", alias.String(), alias.Substitution(alias.Parameters))
	}

	flags := getCustomCmd.Flags()
	flags.StringVar(&getCustomFlags.filterString, "filter", "", "Provide custom LDAP query filter")

	return getCustomCmd
}()

// Runs always prior to "run"
func getCustomPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getCustomPersistentPreRun"})
	logger.Debug("Executing")

	if len(getFlags.searchArguments.Attributes) == 0 {
		getFlags.searchArguments.Attributes.Append(defaultCustomGetAttributes...)
	}
	logger.WithField("searchArguments.Attributes", getFlags.searchArguments.Attributes).Debug("Set")

	apputil.AskFilterString(cmd, "filter", &getCustomFlags.filterString, &getFlags.searchArguments)
	logger.WithField("searchArguments.Filter", getFlags.searchArguments.Filter.String()).Debug("Set")
}
