package commands

import (
	ldap "github.com/go-ldap/ldap/v3"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
	cobra "github.com/spf13/cobra"
)

var editGroupFlags struct {
	id             string
	addMembers     []string
	deleteMembers  []string
	replaceMembers []string
}

var editGroupCmd = func() *cobra.Command {
	editGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Get a group in the directory to edit",
		Long: "Get a group in the directory to edit.\n\n" +
			"In the case no group members are provided, the raw modify request (LDIF) will be requested.\n" +
			"Group members can be provided as lists delimited by semicolon.",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"edit --path \"DC=example,DC=com\" " +
			"group --group-id \"uix12345\"",
		PersistentPreRun: editGroupPersistentPreRun,
		PreRun:           editChildCommandPreRun,
		Run:              editGroupRun,
		PostRun:          editChildCommandPostRun,
	}

	flags := editGroupCmd.Flags()
	flags.StringVar(&editGroupFlags.id, "group-id", "", "Specify group ID (common name, DN or SAN)")
	flags.StringArrayVar(&editGroupFlags.addMembers, "add-member", nil, "Add member with given distinguished name to the group")
	flags.StringArrayVar(&editGroupFlags.deleteMembers, "remove-member", nil, "Remove member with given distinguished name from the group")
	flags.StringArrayVar(&editGroupFlags.replaceMembers, "replace-member", nil, "Replace members of the group with the provided ones identified by their distinguished names")

	return editGroupCmd
}()

func editGroupPersistentPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editGroupPersistentPreRun"})
	logger.Debug("Executing")

	var filters []filter.Filter
	if editGroupFlags.id != "" {
		filters = append(filters, filter.ByID(filter.EscapeFilter(editGroupFlags.id)))
	}

	if len(editGroupFlags.addMembers) > 0 {
		editGroupFlags.addMembers = supererrors.ExceptFn(supererrors.W(libutil.RebuildStringSliceFlag(editGroupFlags.addMembers, ';')))
	}

	if len(editGroupFlags.deleteMembers) > 0 {
		editGroupFlags.deleteMembers = supererrors.ExceptFn(supererrors.W(libutil.RebuildStringSliceFlag(editGroupFlags.deleteMembers, ';')))
	}

	if len(editGroupFlags.replaceMembers) > 0 {
		editGroupFlags.replaceMembers = supererrors.ExceptFn(supererrors.W(libutil.RebuildStringSliceFlag(editGroupFlags.replaceMembers, ';')))
	}

	editFlags.searchArguments.Filter = filter.And(filter.IsGroup(), filters...)
	logger.WithField("searchArguments.Filter", editFlags.searchArguments.Filter.String()).Debug("Set")
}

func editGroupRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editGroupRun"})
	logger.Debug("Executing")

	requests := editFlags.requests
	entry := requests.Entries[0]
	request := ldap.NewModifyRequest(entry.Entry.DN, nil)

	var changes int
	if len(editGroupFlags.replaceMembers) > 0 {
		changes += len(editGroupFlags.replaceMembers)
		request.Delete(attributes.Members().String(), editGroupFlags.replaceMembers)
		logger.WithField("flag", "replace-member").Debug("Set")
	}

	if len(editGroupFlags.addMembers) > 0 {
		changes += len(editGroupFlags.addMembers)
		request.Add(attributes.Members().String(), editGroupFlags.addMembers)
		logger.WithField("flag", "add-member").Debug("Set")
	}

	if len(editGroupFlags.deleteMembers) > 0 {
		changes += len(editGroupFlags.deleteMembers)
		request.Delete(attributes.Members().String(), editGroupFlags.deleteMembers)
		logger.WithField("flag", "delete-member").Debug("Set")
	}

	if changes == 0 {
		supererrors.Except(apputil.AskLDAPDataInterchangeFormat(requests, editFlags.editor))
		editFlags.requests = requests
		logger.WithField("editor", editFlags.editor).Debug("Asked")

		return
	}

	entry.Modify = request
	requests.Entries[0] = entry
	editFlags.requests = requests
}
