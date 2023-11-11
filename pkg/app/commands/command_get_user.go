package commands

import (
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

var getUserFlags struct {
	id          string
	enabled     bool
	expired     bool
	memberOf    string
	recursively bool
}

var getUserCmd = func() *cobra.Command {
	getUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Get a user(s) in the directory",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"get --path \"DC=example,DC=com\" --select \"accountExpires,sAmAccountName\" " +
			"user --user-id \"uix12345\" --enabled",
		PersistentPreRun: getUserPersistentPreRun,
		Run:              getXRun,
	}

	flags := getUserCmd.Flags()
	flags.StringVar(&getUserFlags.id, "user-id", "", "Specify user ID (common name, DN, SAN or UPN)")
	flags.BoolVar(&getUserFlags.enabled, "enabled", false, "Search explicitly for enabled users")
	flags.BoolVar(&getUserFlags.expired, "expired", false, "Search explicitly for expired users")
	flags.StringVar(&getUserFlags.memberOf, "member-of", "", "Search users being member of given group")
	flags.BoolVar(&getUserFlags.recursively, "recursively", false, "Consider nested group membership")

	return getUserCmd
}()

func getUserPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	getFlags.searchArguments.Attributes = append(getFlags.searchArguments.Attributes, defaultUserGetAttributes...)

	var filters []filter.Filter
	if getUserFlags.id != "" {
		filters = append(filters, filter.Or(
			filter.Filter{Attribute: attributes.SamAccountName(), Value: getUserFlags.id},
			filter.Filter{Attribute: attributes.UserPrincipalName(), Value: getUserFlags.id},
			filter.Filter{Attribute: attributes.Name(), Value: getUserFlags.id},
			filter.Filter{Attribute: attributes.DistinguishedName(), Value: getUserFlags.id},
		))
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

	switch {
	case getUserFlags.memberOf != "" && getUserFlags.recursively:
		filters = append(filters, filter.Filter{Attribute: attributes.MemberOf(), Value: getUserFlags.memberOf, Rule: attributes.LDAP_MATCHING_RULE_IN_CHAIN})

	case getUserFlags.memberOf != "" && !getUserFlags.recursively:
		filters = append(filters, filter.Filter{Attribute: attributes.MemberOf(), Value: getUserFlags.memberOf})

	}

	getFlags.searchArguments.Filter = filter.And(filter.IsUser(), filters...)
}
