package commands

import (
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var defaultUserGetAttributes = attributes.Attributes{
	attributes.CommonName(),
	attributes.DisplayName(),
	attributes.DistinguishedName(),
	attributes.Name(),
	attributes.Mail(),
	attributes.ObjectCategory(),
	attributes.ObjectClass(),
	attributes.SamAccountName(),
	attributes.SamAccountType(),
	attributes.UserAccountControl(),
	attributes.UserPrincipalName(),
}

var getUserFlags = &struct {
	id          string
	enabled     bool
	expired     bool
	memberOf    string
	recursively bool
}{}

var getUserCmd = func() *cobra.Command {
	getUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Get a user(s) in the directory",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"get --path \"DC=example,DC=com\" --select \"accountExpires,sAmAccountName\" " +
			"user --user-id \"uix12345\" --enabled",
		PersistentPreRun: getUserPersistentPreRun,
		Run:              getChildCommandRun,
	}

	flags := getUserCmd.Flags()
	flags.StringVar(&getUserFlags.id, "user-id", "", "Specify user ID (common name, DN, SAN or UPN)")
	flags.BoolVar(&getUserFlags.enabled, "enabled", false, "Search explicitly for enabled users")
	flags.BoolVar(&getUserFlags.expired, "expired", false, "Search explicitly for expired users")
	flags.StringVar(&getUserFlags.memberOf, "member-of", "", "Search users being member of given group")
	if getUserFlags.memberOf != "" {
		flags.BoolVar(&getUserFlags.recursively, "recursively", false, "Consider nested group membership")
	}

	return getUserCmd
}()

func getUserPersistentPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getUserPersistentPreRun"})
	logger.Debug("Executing")

	getFlags.searchArguments.Attributes.Append(defaultUserGetAttributes...)
	logger.WithField("searchArguments.Attributes", getFlags.searchArguments.Attributes).Debug("Set")

	var filters []filter.Filter
	if getUserFlags.id != "" {
		filters = append(filters, filter.ByID(getUserFlags.id))
	}

	switch wasProvided := cmd.Flags().Changed("enabled"); {
	case wasProvided && getUserFlags.enabled:
		filters = append(filters, filter.IsEnabled())

	case wasProvided && !getUserFlags.enabled:
		filters = append(filters, filter.Not(filter.IsEnabled()))

	}

	switch wasProvided := cmd.Flags().Changed("expired"); {
	case wasProvided && getUserFlags.expired:
		filters = append(filters, filter.Not(filter.HasNotExpired(true)))

	case wasProvided && !getUserFlags.expired:
		filters = append(filters, filter.HasNotExpired(false))

	}

	if getUserFlags.memberOf != "" {
		filters = append(filters, filter.MemberOf(getUserFlags.memberOf, getUserFlags.recursively))
	}

	getFlags.searchArguments.Filter = filter.And(filter.IsUser(), filters...)
	logger.WithField("searchArguments.Filter", getFlags.searchArguments.Filter.String()).Debug("Set")
}
