package util

import (
	supererrors "github.com/sarumaj/go-super/errors"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

func AskFilterString(cmd *cobra.Command, flagName string, filterString *string, searchArguments *client.SearchArguments) {
	if *filterString == "" {
		var args []string
		_ = supererrors.ExceptFn(supererrors.W(AskString(cmd, flagName, &args, false, "")))
		supererrors.Except(cmd.ParseFlags(args))
	}

	searchArguments.Filter = *supererrors.ExceptFn(supererrors.W(filter.ParseRaw(*filterString)))
}

func AskID(cmd *cobra.Command, flagName string, id *string, searchArguments *client.SearchArguments) {
	if *id == "" {
		var args []string
		_ = supererrors.ExceptFn(supererrors.W(AskString(cmd, flagName, &args, false, "")))
		supererrors.Except(cmd.ParseFlags(args))

		searchArguments.Filter = filter.ByID(*id)
	}
}
