package commands

import (
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
	cobra "github.com/spf13/cobra"
)

// Command options
var editGroupFlags struct {
	id               string   `flag:"group-id"`
	addMembers       []string `flag:"add-member"`
	deleteMembers    []string `flag:"remove-member"`
	replaceMembers   []string `flag:"replace-member"`
	membersAttribute string   `flag:"member-attribute"`
}

// "group" command
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
	flags.StringVar(&editGroupFlags.membersAttribute, "member-attribute", attributes.Members().String(), "Configure custom attribute name for variant directory schema")

	return editGroupCmd
}()

// Runs always prior to "run"
func editGroupPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editGroupPersistentPreRun"})
	logger.Trace("Executing")

	apputil.AskID(cmd, "group-id", &editGroupFlags.id, &editFlags.searchArguments)

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
	logger.WithField("searchArguments.Filter", editFlags.searchArguments.Filter.String()).Trace("Set")
}

// Actual "run" prepares a modify request
func editGroupRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editGroupRun"})
	logger.Trace("Executing")

	requests := editFlags.requests
	entry := requests.Entries[0]

	request := client.ModifyGroupMembersRequest(
		entry.Entry.DN,
		editGroupFlags.addMembers,
		editGroupFlags.deleteMembers,
		editGroupFlags.replaceMembers,
		attributes.Attribute{LDAPDisplayName: editGroupFlags.membersAttribute},
	)

	if len(request.Changes) == 0 {
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskLDAPDataInterchangeFormat(requests, editFlags.editor)))
		editFlags.requests = requests
		logger.WithField("editor", editFlags.editor).Trace("Asked")

		return
	}

	entry.Modify = request
	requests.Entries[0] = entry
	editFlags.requests = requests
	logger.WithField("flag", "modify-member").Trace("Set")
}
