package commands

import (
	"fmt"
	"runtime"

	color "github.com/fatih/color"
	ldif "github.com/go-ldap/ldif"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/v2/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/v2/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/v2/pkg/lib/client"
	attributes "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/attributes"
	progressbar "github.com/schollz/progressbar/v3"
	cobra "github.com/spf13/cobra"
)

// Command options
var editFlags struct {
	editor          string `flag:"editor"`
	requests        *ldif.LDIF
	searchArguments client.SearchArguments
}

// "edit" command
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

// Runs prior to "run" and executes search query for a child command
func editChildCommandPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editChildCommandPreRun"})
	logger.Trace("Executing")

	logger.WithFields(apputil.GetFieldsForBind(&rootFlags.bindParameters, &rootFlags.dialOptions)).Trace("Connecting")
	conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
		&rootFlags.bindParameters,
		&rootFlags.dialOptions,
	)))

	logger.WithFields(apputil.GetFieldsForSearch(&editFlags.searchArguments)).Trace("Querying")
	_, requests := supererrors.ExceptFn2(supererrors.W2(client.Search(conn, editFlags.searchArguments, progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(apputil.Stdout()),
		progressbar.OptionEnableColorCodes(apputil.IsColorEnabled()),
	))))

	if len(requests.Entries) == 0 {
		apputil.PrintlnAndExit(1, "There is nothing to edit matching: %q", editFlags.searchArguments.Filter)
	}

	if len(requests.Entries) > 1 {
		apputil.PrintlnAndExit(1, "Provided filter led to a ubiquitous result set (%d entries for %q)", len(requests.Entries), editFlags.searchArguments.Filter)
	}

	editFlags.requests = requests
}

// Runs after a "run" and executes a modify request for a child command
func editChildCommandPostRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editChildCommandPostRun"})
	logger.Trace("Executing")

	logger.WithFields(apputil.GetFieldsForBind(&rootFlags.bindParameters, &rootFlags.dialOptions)).Trace("Connecting")
	conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
		&rootFlags.bindParameters,
		&rootFlags.dialOptions,
	)))

	for i, entry := range editFlags.requests.Entries {
		if entry.Add != nil {
			logger.WithField(fmt.Sprintf("requests.Entries[%d].Add", i), entry.Add).Trace("Creating")
			supererrors.Except(conn.Add(entry.Add))
		}

		if entry.Del != nil {
			logger.WithField(fmt.Sprintf("requests.Entries[%d].Del", i), entry.Del).Trace("Applying")
			supererrors.Except(conn.Del(entry.Del))
		}

		if entry.Modify != nil {
			logger.WithField(fmt.Sprintf("requests.Entries[%d].Modify", i), entry.Modify).Trace("Applying")
			supererrors.Except(conn.Modify(entry.Modify))
		}
	}

	apputil.PrintlnAndExit(0, apputil.PrintColors(color.HiGreenString, "Successfully applied modifications"))
}

// Runs prior to "run" and sets search query options (inherited by child commands)
func editPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editPersistentPreRun"})
	logger.Trace("Executing")

	if editFlags.searchArguments.Path == "" {
		editFlags.searchArguments.Path = rootFlags.dialOptions.URL.ToBaseDirectoryPath()
		logger.WithField("searchArguments.Path", editFlags.searchArguments.Path).Trace("Set")
	}

	// select all properties, even the unregistered ones
	editFlags.searchArguments.Attributes = attributes.LookupMany(false, "*")
	logger.WithField("searchArguments.attributes", editFlags.searchArguments.Attributes).Trace("Set")
}

// Runs in an interactive mode by asking user to provide options for the command
func editRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "editRun"})
	logger.Trace("Executing")

	child := supererrors.ExceptFn(supererrors.W(apputil.AskCommand(cmd, editGroupCmd)))
	logger = logger.WithField("child", child.Name())

	var args []string
	switch child {

	case editCustomCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "filter", &args, false, "")))
		logger.WithFields(apputil.Fields{"flag": "filter", "args": args}).Trace("Asked")

	case editGroupCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "group-id", &args, false, "")))

		switch {
		case
			supererrors.ExceptFn(supererrors.W(apputil.AskMultiline(child, "add-member", &args))),
			supererrors.ExceptFn(supererrors.W(apputil.AskMultiline(child, "remove-member", &args))),
			supererrors.ExceptFn(supererrors.W(apputil.AskMultiline(child, "replace-member", &args))):

			supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "member-attribute", &args, false, attributes.Members().String())))
		}

		logger.WithFields(apputil.Fields{"flags": []string{"group-id", "add-member", "remove-member", "replace-member"}, "args": args}).Trace("Asked")

	case editUserCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "user-id", &args, false, "")))
		if supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "new-password", &args, true, ""))) {
			_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "password-attribute", &args, false, attributes.UnicodePassword().String())))
		}
		logger.WithFields(apputil.Fields{"flags": []string{"user-id", "new-password", "password-attribute"}, "args": args}).Trace("Asked")

	}

	_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "path", &args, false, rootFlags.dialOptions.URL.ToBaseDirectoryPath())))
	logger.WithFields(apputil.Fields{"flag": "path", "args": args}).Trace("Asked")

	supererrors.Except(child.ParseFlags(args))
	logger.Trace("Parsed")

	child.PersistentPreRun(child, args)
	child.PreRun(child, args)
	child.Run(child, args)
	child.PostRun(child, args)
}
