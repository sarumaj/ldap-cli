package commands

import (
	"fmt"
	"runtime"

	ldif "github.com/go-ldap/ldif"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	cobra "github.com/spf13/cobra"
)

var editFlags = &struct {
	editor          string
	requests        *ldif.LDIF
	searchArguments client.SearchArguments
}{}

var editCmd = func() *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a directory object",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"edit --path \"DC=example,DC=com\" <command>",
		PersistentPreRun: editPersistentPreRun,
		Run:              editRun,
	}

	editor := "vi"
	if runtime.GOOS == "windows" {
		editor = "notepad"
	}

	flags := editCmd.PersistentFlags()
	flags.StringVar(&editFlags.editor, "editor", editor, "Specify editor to modify *.ldif files")
	flags.StringVar(&editFlags.searchArguments.Path, "path", "", "Specify the query path to search the directory objects in (per default path is derivated from the address of domain controller)")

	editCmd.AddCommand(editCustomCmd, editGroupCmd, editUserCmd)

	return editCmd
}()

func editChildCommandPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editChildCommandPreRun"})
	logger.Debug("Executing")

	logger.WithFields(apputil.GetFieldsForBind(&rootFlags.bindParameters, &rootFlags.dialOptions)).Debug("Connecting")
	conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
		&rootFlags.bindParameters,
		&rootFlags.dialOptions,
	)))

	logger.WithFields(apputil.GetFieldsForSearch(&editFlags.searchArguments)).Debug("Querying")
	_, requests := supererrors.ExceptFn2(supererrors.W2(client.Search(conn, editFlags.searchArguments)))

	if len(requests.Entries) == 0 {
		apputil.PrintlnAndExit("There is nothing to edit matching: %q", editFlags.searchArguments.Filter)
	}

	if len(requests.Entries) > 1 {
		apputil.PrintlnAndExit("Provided filter led to a ubiquitous result set (%d entries for %q)", len(requests.Entries), editFlags.searchArguments.Filter)
	}

	editFlags.requests = requests
}

func editChildCommandPostRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editChildCommandPostRun"})
	logger.Debug("Executing")

	logger.WithFields(apputil.GetFieldsForBind(&rootFlags.bindParameters, &rootFlags.dialOptions)).Debug("Connecting")
	conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
		&rootFlags.bindParameters,
		&rootFlags.dialOptions,
	)))

	for i, entry := range editFlags.requests.Entries {
		if entry.Add != nil {
			logger.WithField(fmt.Sprintf("requests.Entries[%d].Add", i), entry.Add).Debug("Creating")
			supererrors.Except(conn.Add(entry.Add))
		}

		if entry.Del != nil {
			logger.WithField(fmt.Sprintf("requests.Entries[%d].Del", i), entry.Del).Debug("Applying")
			supererrors.Except(conn.Del(entry.Del))
		}

		if entry.Modify != nil {
			logger.WithField(fmt.Sprintf("requests.Entries[%d].Modify", i), entry.Modify).Debug("Applying")
			supererrors.Except(conn.Modify(entry.Modify))
		}
	}
}

func editPersistentPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editPersistentPreRun"})
	logger.Debug("Executing")

	if editFlags.searchArguments.Path == "" {
		editFlags.searchArguments.Path = rootFlags.dialOptions.URL.ToBaseDirectoryPath()
		logger.WithField("searchArguments.Path", editFlags.searchArguments.Path).Debug("Set")
	}

	// select all properties, even the unregistered ones
	editFlags.searchArguments.Attributes = attributes.LookupMany(false, "*")
	logger.WithField("searchArguments.attributes", editFlags.searchArguments.Attributes).Debug("Set")
}

func editRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editRun"})
	logger.Debug("Executing")

	child := supererrors.ExceptFn(supererrors.W(apputil.AskCommand(cmd, editGroupCmd)))
	logger = logger.WithField("child", child.Name())

	var args []string
	switch child {

	case editCustomCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "filter", &args, false)))
		logger.WithFields(apputil.Fields{"flag": "filter", "args": args}).Debug("Asked")

	case editGroupCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "group-id", &args, false)))

		switch {
		case
			supererrors.ExceptFn(supererrors.W(apputil.AskMultiline(child, "add-member", &args))),
			supererrors.ExceptFn(supererrors.W(apputil.AskMultiline(child, "remove-member", &args))),
			supererrors.ExceptFn(supererrors.W(apputil.AskMultiline(child, "replace-member", &args))):

			supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "member-attribute", &args, false)))
		}

		logger.WithFields(apputil.Fields{"flags": []string{"group-id", "add-member", "remove-member", "replace-member"}, "args": args}).Debug("Asked")

	case editUserCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "user-id", &args, false)))
		if supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "password", &args, true))) {
			_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "password-attribute", &args, false)))
		}
		logger.WithFields(apputil.Fields{"flags": []string{"user-id", "password", "password-attribute"}, "args": args}).Debug("Asked")

	}

	_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "path", &args, false)))
	logger.WithFields(apputil.Fields{"flag": "path", "args": args}).Debug("Asked")

	supererrors.Except(child.ParseFlags(args))
	logger.Debug("Parsed")

	// since the flags could have changed, pre run must be invoked again
	cmd.PersistentPreRun(cmd, nil)
	child.PersistentPreRun(child, nil)
	child.PreRun(child, nil)
	child.Run(child, nil)
	child.PostRun(child, nil)
}
