package commands

import (
	ldap "github.com/go-ldap/ldap/v3"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var editUserFlags struct {
	id       string
	password string
}

var editUserCmd = func() *cobra.Command {
	editUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Get a user in the directory to edit",
		Long: "Get a user in the directory to edit.\n\n" +
			"In the case no password is provided, the raw modify request (LDIF) will be requested.",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"edit --path \"DC=example,DC=com\" --select \"accountExpires,sAmAccountName\" " +
			"user --user-id \"uix12345\"",
		PersistentPreRun: editUserPersistentPreRun,
		PreRun:           editChildCommandPreRun,
		Run:              editUserRun,
		PostRun:          editChildCommandPostRun,
	}

	flags := editUserCmd.Flags()
	flags.StringVar(&editUserFlags.id, "user-id", "", "Specify user ID (common name, DN, SAN or UPN)")
	flags.StringVar(&editUserFlags.password, "password", "", "Provide new password for the user to change (leave empty to keep current)")

	return editUserCmd
}()

func editUserPersistentPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editUserPersistentPreRun"})
	logger.Debug("Executing")

	var filters []filter.Filter
	if editUserFlags.id != "" {
		filters = append(filters, filter.ByID(filter.EscapeFilter(editUserFlags.id)))
	}

	editFlags.searchArguments.Filter = filter.And(filter.IsUser(), filters...)
	logger.WithField("searchArguments.Filter", editFlags.searchArguments.Filter.String()).Debug("Set")
}

func editUserRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editUserRun"})
	logger.Debug("Executing")

	requests := editFlags.requests
	entry := requests.Entries[0]

	if editUserFlags.password == "" {
		supererrors.Except(apputil.AskLDAPDataInterchangeFormat(requests, editFlags.editor))
		editFlags.requests = requests
		logger.WithField("editor", editFlags.editor).Debug("Asked")

		return
	}

	request := ldap.NewModifyRequest(entry.Entry.DN, nil)
	request.Replace(attributes.UserPassword().String(), []string{editUserFlags.password})

	entry.Modify = request
	requests.Entries[0] = entry
	editFlags.requests = requests

	logger.WithField("flag", "password").Debug("Set")
}
