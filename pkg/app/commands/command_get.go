package commands

import (
	"fmt"
	"os"

	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
	progressbar "github.com/schollz/progressbar/v3"
	cobra "github.com/spf13/cobra"
)

// Default attributes for search queries
var defaultGetAttributes = map[string]attributes.Attributes{
	getCustomCmd.Name(): defaultCustomGetAttributes,
	getGroupCmd.Name():  defaultGroupGetAttributes,
	getUserCmd.Name():   defaultUserGetAttributes,
}

// Command options
var getFlags struct {
	format           string `flag:"format"`
	searchArguments  client.SearchArguments
	selectAttributes []string `flag:"select"`
	output           string
}

// "get" command
var getCmd = func() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a directory object",
		Example: "ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"get --path \"DC=example,DC=com\" --select \"accountExpires,sAmAccountName\" <command>\n" +
			"ldap-cli --user \"DOMAIN\\\\user\" --password \"password\" --url \"ldaps://example.com:636\" " +
			"get --path \"DC=example,DC=com\" --select \"*\" <command>\n",
		PersistentPreRun: getPersistentPreRun,
		Run:              getRun,
	}

	flags := getCmd.PersistentFlags()
	flags.StringVar(&getFlags.format, "format", "default", fmt.Sprintf("Output format (supported: [%v])", apputil.ListSupportedFormats(true)))
	flags.StringVar(&getFlags.searchArguments.Path, "path", "", "Specify the query path to search the directory objects in (per default path is derivated from the address of domain controller)")
	flags.StringArrayVar(&getFlags.selectAttributes, "select", []string{}, "Select specific object attributes (if not provided default attributes are being selected)")
	flags.StringVar(&getFlags.output, "output", "stdout", "Output to file")

	getCmd.AddCommand(getCustomCmd, getGroupCmd, getUserCmd)

	return getCmd
}()

// Runner for a sub-command of the "get" command.
// It will bind to a domain controller and execute the search query
func getChildCommandRun(cmd *cobra.Command, args []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getChildCommandRun"})
	logger.Trace("Executing")

	logger.WithFields(apputil.GetFieldsForBind(&rootFlags.bindParameters, &rootFlags.dialOptions)).Trace("Connecting")
	conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
		&rootFlags.bindParameters,
		&rootFlags.dialOptions,
	)))

	logger.WithFields(apputil.GetFieldsForSearch(&getFlags.searchArguments)).Trace("Querying")
	results, requests := supererrors.ExceptFn2(supererrors.W2(client.Search(conn, getFlags.searchArguments, progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(apputil.Stdout()),
		progressbar.OptionEnableColorCodes(apputil.IsColorEnabled()),
	))))

	logger.WithField("format", getFlags.format).WithField("output", getFlags.output).Trace("Rendering")
	if getFlags.output == "stdout" {
		supererrors.Except(apputil.Flush(results, requests, getFlags.format, apputil.Stdout()))
	} else {
		f := supererrors.ExceptFn(supererrors.W(os.OpenFile(getFlags.output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)))
		supererrors.Except(apputil.Flush(results, requests, apputil.SniffFormat(f.Name(), getFlags.format), f))
		supererrors.Except(f.Close(), os.ErrClosed)
	}

}

// Runs always prior to "run" (inherited by child commands of the "get" command).
// Search request parameters will be set from command options
func getPersistentPreRun(cmd *cobra.Command, args []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, args)

	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getPersistentPreRun"})
	logger.Trace("Executing")

	if getFlags.searchArguments.Path == "" {
		getFlags.searchArguments.Path = rootFlags.dialOptions.URL.ToBaseDirectoryPath()
		logger.WithField("searchArguments.Path", getFlags.searchArguments.Path).Trace("Set")
	}

	if len(getFlags.selectAttributes) > 0 {
		selected := supererrors.ExceptFn(supererrors.W(libutil.RebuildStringSliceFlag(getFlags.selectAttributes, ',')))
		getFlags.searchArguments.Attributes = attributes.LookupMany(false, selected...)
		logger.WithField("searchArguments.Attributes", getFlags.searchArguments.Attributes).Trace("Set")
	}
}

// Runs "get" command in interactive mode by asking user to provide values for command parameters
func getRun(cmd *cobra.Command, args []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getRun"})
	logger.Trace("Executing")

	child := supererrors.ExceptFn(supererrors.W(apputil.AskCommand(cmd, getUserCmd)))
	logger = logger.WithField("child", child.Name())

	switch child {

	case getCustomCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "filter", &args, false, "")))
		logger.WithFields(apputil.Fields{"flag": "filter", "args": args}).Trace("Asked")

	case getGroupCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "group-id", &args, false, "")))
		logger.WithFields(apputil.Fields{"flag": "group-id", "args": args}).Trace("Asked")

	case getUserCmd:
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "user-id", &args, false, "")))
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskBool(child, "enabled", &args)))
		_ = supererrors.ExceptFn(supererrors.W(apputil.AskBool(child, "expired", &args)))
		if supererrors.ExceptFn(supererrors.W(apputil.AskMultiline(child, "member-of", &args))) {
			_ = supererrors.ExceptFn(supererrors.W(apputil.AskBool(child, "recursively", &args)))
		}
		logger.WithFields(apputil.Fields{"flags": []string{"user-id", "enabled", "expired", "member-of", "recursively"}, "args": args}).Trace("Asked")

	}

	options, defaults := append([]string{"*"}, attributes.LookupMany(true, "*").ToAttributeList()...), defaultGetAttributes[child.Name()].ToAttributeList()
	_ = supererrors.ExceptFn(supererrors.W(apputil.AskStrings(child, "select", options, defaults, &args)))
	_ = supererrors.ExceptFn(supererrors.W(apputil.AskString(child, "path", &args, false, rootFlags.dialOptions.URL.ToBaseDirectoryPath())))
	_ = supererrors.ExceptFn(supererrors.W(apputil.AskStrings(child, "format", []string{"csv", "default", "ldif", "yaml"}, []string{"default"}, &args)))
	logger.WithFields(apputil.Fields{"flags": []string{"select", "path", "format"}, "args": args}).Trace("Asked")

	supererrors.Except(child.ParseFlags(args))
	logger.Trace("Parsed")

	child.PersistentPreRun(child, args)
	child.Run(child, args)
}
