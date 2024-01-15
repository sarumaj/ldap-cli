package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
	cobra "github.com/spf13/cobra"
)

// Default attributes for search query
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

// Command options
var getUserFlags struct {
	id          string   `flag:"user-id"`
	enabled     bool     `flag:"enabled"`
	expired     bool     `flag:"expired"`
	memberOf    []string `flag:"member-of"`
	recursively bool     `flag:"recursively"`
}

// "user" command
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
	flags.StringArrayVar(&getUserFlags.memberOf, "member-of", nil, "Search users being member of given group")
	flags.BoolVar(&getUserFlags.recursively, "recursively", false, "Consider nested group membership")

	return getUserCmd
}()

// Runs always prior to "run"
func getUserPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getUserPersistentPreRun"})
	logger.Debug("Executing")

	apputil.AskID(cmd, "user-id", &getUserFlags.id, &getFlags.searchArguments)

	if len(getFlags.searchArguments.Attributes) == 0 {
		getFlags.searchArguments.Attributes.Append(defaultUserGetAttributes...)
	}
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
		filters = append(filters, filter.HasExpired())

	case wasProvided && !getUserFlags.expired:
		filters = append(filters, filter.Not(filter.HasExpired()))

	}

	if len(getUserFlags.memberOf) > 0 {
		getUserFlags.memberOf = supererrors.ExceptFn(supererrors.W(libutil.RebuildStringSliceFlag(getUserFlags.memberOf, ';')))
	}

	for _, memberOf := range getUserFlags.memberOf {
		filters = append(filters, filter.MemberOf(memberOf, getUserFlags.recursively))
	}

	getFlags.searchArguments.Filter = filter.And(filter.IsUser(), filters...)
	logger.WithField("searchArguments.Filter", getFlags.searchArguments.Filter.String()).Debug("Set")
}
