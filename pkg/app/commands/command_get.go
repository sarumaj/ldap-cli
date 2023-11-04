package commands

import (
	"bytes"
	"encoding/csv"
	"fmt"

	color "github.com/fatih/color"
	ldif "github.com/go-ldap/ldif"
	yaml "github.com/goccy/go-yaml"
	lexer "github.com/goccy/go-yaml/lexer"
	printer "github.com/goccy/go-yaml/printer"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	libutil "github.com/sarumaj/ldap-cli/pkg/lib/util"
	cobra "github.com/spf13/cobra"
)

var defaultGetAttributes = map[string]attributes.Attributes{
	getCustomCmd.Name(): defaultCustomGetAttributes,
	getGroupCmd.Name():  defaultGroupGetAttributes,
	getUserCmd.Name():   defaultUserGetAttributes,
}

var getFlags = &struct {
	format           string
	searchArguments  client.SearchArguments
	selectAttributes []string
}{}

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

	getCmd.AddCommand(getCustomCmd, getGroupCmd, getUserCmd)

	return getCmd
}()

func getChildCommandRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getChildCommandRun"})
	logger.Debug("Executing")

	logger.WithFields(apputil.GetFieldsForBind(&rootFlags.bindParameters, &rootFlags.dialOptions)).Debug("Connecting")
	conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
		&rootFlags.bindParameters,
		&rootFlags.dialOptions,
	)))

	logger.WithFields(apputil.GetFieldsForSearch(&getFlags.searchArguments)).Debug("Querying")
	results, requests := supererrors.ExceptFn2(supererrors.W2(client.Search(conn, getFlags.searchArguments)))

	logger.WithField("format", getFlags.format).Debug("Rendering")
	switch getFlags.format {

	case apputil.CSV:
		lines := make([][]string, len(results)+1)
		for i, m := range results {
			for _, a := range attributes.Map(m).Keys() {
				if i == 0 {
					lines[0] = append(lines[0], a.String())
				}

				lines[i+1] = append(lines[i+1], fmt.Sprint(m[a]))
			}
		}
		supererrors.Except(csv.NewWriter(apputil.Stdout()).WriteAll(lines))

	case apputil.LDIF:
		data := supererrors.ExceptFn(supererrors.W(ldif.Marshal(requests)))
		_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), data)))

	case apputil.YAML:
		if len(results) == 1 {
			supererrors.Except(yaml.NewEncoder(apputil.Stdout(), yaml.Indent(2)).Encode(results[0]))
		} else {
			supererrors.Except(yaml.NewEncoder(apputil.Stdout(), yaml.Indent(2)).Encode(map[string]any{"Results": results}))
		}

	default:
		buffer := bytes.NewBuffer(nil)
		if len(results) == 1 {
			supererrors.Except(yaml.NewEncoder(buffer, yaml.Indent(2)).Encode(results[0]))
		} else {
			supererrors.Except(yaml.NewEncoder(buffer, yaml.Indent(2)).Encode(map[string]any{"Results": results}))
		}

		tokens := lexer.Tokenize(buffer.String())
		buffer.Reset()

		_ = supererrors.ExceptFn(supererrors.W(fmt.Fprintln(apputil.Stdout(), (&printer.Printer{
			Bool: func() *printer.Property {
				return &printer.Property{Prefix: fmt.Sprintf("\x1b[%dm", color.FgHiMagenta), Suffix: "\x1b[0m"}
			},
			Number: func() *printer.Property {
				return &printer.Property{Prefix: fmt.Sprintf("\x1b[%dm", color.FgHiMagenta), Suffix: "\x1b[0m"}
			},
			MapKey: func() *printer.Property {
				return &printer.Property{Prefix: fmt.Sprintf("\x1b[%dm", color.FgHiCyan), Suffix: "\x1b[0m"}
			},
			String: func() *printer.Property {
				return &printer.Property{Prefix: fmt.Sprintf("\x1b[%dm", color.FgHiGreen), Suffix: "\x1b[0m"}
			},
		}).PrintTokens(tokens))))
	}
}

func getPersistentPreRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getPersistentPreRun"})
	logger.Debug("Executing")

	if getFlags.searchArguments.Path == "" {
		getFlags.searchArguments.Path = rootFlags.dialOptions.URL.ToBaseDirectoryPath()
		logger.WithField("searchArguments.Path", getFlags.searchArguments.Path).Debug("Set")
	}

	if len(getFlags.selectAttributes) > 0 {
		selected := supererrors.ExceptFn(supererrors.W(libutil.RebuildStringSliceFlag(getFlags.selectAttributes, ',')))
		getFlags.searchArguments.Attributes = attributes.LookupMany(false, selected...)
		logger.WithField("searchArguments.Attributes", getFlags.searchArguments.Attributes).Debug("Set")
	}
}

func getRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithFields(apputil.Fields{"command": cmd.CommandPath(), "step": "getRun"})
	logger.Debug("Executing")

	child := supererrors.ExceptFn(supererrors.W(apputil.AskCommand(cmd, getUserCmd)))
	logger = logger.WithField("child", child.Name())

	var args []string
	switch child {

	case getCustomCmd:
		supererrors.Except(apputil.AskString(child, "filter", &args, false))
		logger.WithFields(apputil.Fields{"flag": "filter", "args": args}).Debug("Asked")

	case getGroupCmd:
		supererrors.Except(apputil.AskString(child, "group-id", &args, false))
		logger.WithFields(apputil.Fields{"flag": "group-id", "args": args}).Debug("Asked")

	case getUserCmd:
		supererrors.Except(apputil.AskString(child, "user-id", &args, false))
		supererrors.Except(apputil.AskBool(child, "enabled", &args))
		supererrors.Except(apputil.AskBool(child, "expired", &args))
		supererrors.Except(apputil.AskString(child, "member-of", &args, false))
		supererrors.Except(apputil.AskBool(child, "recursively", &args))
		logger.WithFields(apputil.Fields{"flags": []string{"user-id", "enabled", "expired", "member-of", "recursively"}, "args": args}).Debug("Asked")

	}

	options, defaults := append([]string{"*"}, attributes.LookupMany(true, "*").ToAttributeList()...), defaultGetAttributes[child.Name()].ToAttributeList()
	supererrors.Except(apputil.AskStrings(child, "select", options, defaults, &args))
	supererrors.Except(apputil.AskString(child, "path", &args, false))
	supererrors.Except(apputil.AskStrings(child, "format", []string{"csv", "default", "ldif", "yaml"}, []string{"default"}, &args))
	logger.WithFields(apputil.Fields{"flags": []string{"select", "path", "format"}, "args": args}).Debug("Asked")

	supererrors.Except(child.ParseFlags(args))
	logger.Debug("Parsed")

	// since the flags could have changed, pre run must be invoked again
	cmd.PersistentPreRun(cmd, nil)
	child.PersistentPreRun(child, nil)
	child.Run(child, nil)
}
