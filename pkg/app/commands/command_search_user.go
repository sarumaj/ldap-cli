package commands

import (
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
	pflag "github.com/spf13/pflag"
)

var searchUserCmd = func() *cobra.Command {
	var id string
	var enabled bool
	var expired bool
	var memberOf string
	var recursively bool

	searchUserCmd := &cobra.Command{
		Use:     "user",
		Short:   "Search a user(s) in the directory",
		Example: "ldap-cli search user",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			parent := cmd.Parent()
			parent.PersistentPreRun(parent, nil)

			searchArguments.Attributes = append(
				searchArguments.Attributes,
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
			)

			var filters []filter.Filter
			if id != "" {
				filters = append(filters, filter.Or(
					filter.Filter{Attribute: attributes.SamAccountName(), Value: id},
					filter.Filter{Attribute: attributes.UserPrincipalName(), Value: id},
					filter.Filter{Attribute: attributes.Name(), Value: id},
					filter.Filter{Attribute: attributes.CommonName(), Value: id},
				))
			}

			enabledFlagSet, expiredFlagSet := false, false
			cmd.Flags().Visit(func(f *pflag.Flag) {
				switch f.Name {
				case "enabled":
					enabledFlagSet = true

				case "expired":
					expiredFlagSet = true

				}
			})

			switch {
			case enabledFlagSet && enabled:
				filters = append(filters, filter.IsEnabled())

			case enabledFlagSet && !enabled:
				filters = append(filters, filter.Not(filter.IsEnabled()))

			case expiredFlagSet && expired:
				filters = append(filters, filter.Not(filter.HasNotExpired(true)))

			case expiredFlagSet && !expired:
				filters = append(filters, filter.HasNotExpired(false))

			}

			switch {
			case memberOf != "" && recursively:
				filters = append(filters, filter.Filter{Attribute: attributes.MemberOf(), Value: memberOf, Rule: attributes.LDAP_MATCHING_RULE_IN_CHAIN})

			case memberOf != "" && !recursively:
				filters = append(filters, filter.Filter{Attribute: attributes.MemberOf(), Value: memberOf})

			}

			searchArguments.Filter = filter.And(filter.IsUser(), filters...)
		},
		Run: performSearch,
	}

	flags := searchUserCmd.Flags()
	flags.StringVar(&id, "user-id", "", "Specify user ID (common name, SAN or UPN)")
	flags.BoolVar(&enabled, "enabled", false, "Search explicitly for enabled users")
	flags.BoolVar(&expired, "expired", false, "Search explicitly for expired users")
	flags.StringVar(&memberOf, "member-of", "", "Search users being member of given group")
	flags.BoolVar(&recursively, "recursively", false, "Consider nested group membership")

	return searchUserCmd
}()
