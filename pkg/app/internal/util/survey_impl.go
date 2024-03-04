package util

import (
	supererrors "github.com/sarumaj/go-super/errors"
	client "github.com/sarumaj/ldap-cli/v2/pkg/lib/client"
	filter "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

// AskFilterString asks for a filter string (LDAP syntax according to RFC 4515)
func AskFilterString(cmd *cobra.Command, flagName string, filterString *string, searchArguments *client.SearchArguments) {
	if *filterString == "" {
		var args []string
		_ = supererrors.ExceptFn(supererrors.W(AskString(cmd, flagName, &args, false, "")))
		supererrors.Except(cmd.ParseFlags(args))
	}

	searchArguments.Filter = *supererrors.ExceptFn(supererrors.W(filter.ParseRaw(*filterString)))
}

// AskID asks for an ID (CN, DN, GUID, SAN, UPN, Name or DisplayName)
func AskID(cmd *cobra.Command, flagName string, id *string, searchArguments *client.SearchArguments) {
	if *id == "" {
		var args []string
		_ = supererrors.ExceptFn(supererrors.W(AskString(cmd, flagName, &args, false, "")))
		supererrors.Except(cmd.ParseFlags(args))

		searchArguments.Filter = filter.ByID(*id)
	}
}
