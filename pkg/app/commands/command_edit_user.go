package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	cobra "github.com/spf13/cobra"
)

var editUserFlags struct {
	id           string `flag:"user-id"`
	password     string `flag:"new-password"`
	pwdAttribute string `flag:"password-attribute"`
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
	flags.StringVar(&editUserFlags.password, "new-password", "", "Provide new password for the user to change (leave empty to keep current)")
	flags.StringVar(&editUserFlags.pwdAttribute, "password-attribute", attributes.UnicodePassword().String(), "Configure custom attribute name for variant directory schema")

	return editUserCmd
}()

func editUserPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editUserPersistentPreRun"})
	logger.Debug("Executing")

	if editUserFlags.id == "" {
		var args []string
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(cmd, "user-id", &args, false, "")))
		supererrors.Except(cmd.ParseFlags(args))
		getFlags.searchArguments.Filter = filter.ByID(editUserFlags.id)
		logger.WithField("searchArguments.Filter", editFlags.searchArguments.Filter).Debug("Asked")
	}

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
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskLDAPDataInterchangeFormat(requests, editFlags.editor)))
		editFlags.requests = requests
		logger.WithField("editor", editFlags.editor).Debug("Asked")

		return
	}

	entry.Modify = client.ModifyPasswordRequest(entry.Entry.DN, editUserFlags.password, attributes.Attribute{LDAPDisplayName: editUserFlags.pwdAttribute})
	requests.Entries[0] = entry
	editFlags.requests = requests

	logger.WithField("flag", "new-password").Debug("Set")
}
