package commands

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	color "github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var defaultGetAttributes = map[string]attributes.Attributes{
	"custom": defaultCustomGetAttributes,
	"group":  defaultGroupGetAttributes,
	"user":   defaultUserGetAttributes,
}

var getFlags struct {
	format           string
	searchArguments  client.SearchArguments
	selectAttributes string
}

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
	flags.StringVar(&getFlags.format, "format", "default", "Output format (supported: [\"csv\", \"default\", \"yaml\"])")
	flags.StringVar(&getFlags.searchArguments.Path, "path", "", "Specify the query path to search the directory objects in (per default path is derivated from the address of domain controller)")
	flags.StringVar(&getFlags.selectAttributes, "select", "", "Select specific object attributes (if not provided default attributes are being selected)")

	getCmd.AddCommand(getCustomCmd, getGroupCmd, getUserCmd)

	return getCmd
}()

func getPersistentPreRun(cmd *cobra.Command, _ []string) {
	parent := cmd.Parent()
	parent.PersistentPreRun(parent, nil)

	if getFlags.searchArguments.Path == "" {
		var components []string
		for _, dc := range strings.Split(dialOptions.URL.Host, ".") {
			if dc == "" {
				continue
			}

			components = append(components, "DC="+dc)
		}
		getFlags.searchArguments.Path = strings.Join(components, ",")
	}

	if len(getFlags.selectAttributes) > 0 {
		reader := csv.NewReader(strings.NewReader(getFlags.selectAttributes))
		reader.LazyQuotes = true
		reader.TrimLeadingSpace = true

		getFlags.searchArguments.Attributes = attributes.LookupMany(supererrors.ExceptFn(supererrors.W(reader.Read()))...)
	}
}

func getRun(cmd *cobra.Command, args []string) {
	var x string
	supererrors.Except(survey.AskOne(&survey.Select{
		Message: "What do you want to search for?",
		Options: []string{getCustomCmd.Name(), getGroupCmd.Name(), getUserCmd.Name()},
	}, &x))

	var child *cobra.Command
	switch x {

	case getCustomCmd.Name():
		child = getCustomCmd
		supererrors.Except(apputil.AskString(child, "filter", &args))

	case getGroupCmd.Name():
		child = getGroupCmd
		supererrors.Except(apputil.AskString(child, "group-id", &args))

	case getUserCmd.Name():
		child = getUserCmd
		supererrors.Except(apputil.AskString(child, "user-id", &args))
		supererrors.Except(apputil.AskBool(child, "enabled", &args))
		supererrors.Except(apputil.AskBool(child, "expired", &args))
		supererrors.Except(apputil.AskString(child, "member-of", &args))
		supererrors.Except(apputil.AskBool(child, "recursively", &args))

	default:
		return

	}

	supererrors.Except(apputil.AskStrings(child, "select", attributes.LookupMany("*").ToAttributeList(), defaultGetAttributes[child.Name()].ToAttributeList(), &args))
	supererrors.Except(apputil.AskStrings(child, "format", []string{"csv", "default", "yaml"}, []string{"default"}, &args))

	supererrors.Except(child.ParseFlags(args))
	child.PersistentPreRun(child, nil)

	child.Run(child, nil)
}

func getXRun(cmd *cobra.Command, _ []string) {
	logger := apputil.Logger.WithField("command", cmd.CommandPath())

	logger.WithFields(logrus.Fields{
		"bindParameters.AuthType":         bindParameters.AuthType.String(),
		"bindParameters.Domain":           bindParameters.Domain,
		"bindParameters.User":             bindParameters.User,
		"bindParameters.PasswordProvided": len(bindParameters.Password) > 0,
		"dialOptions.MaxRetries":          dialOptions.MaxRetries,
		"dialOptions.SizeLimit":           dialOptions.SizeLimit,
		"dialOptions.TLSEnabled":          dialOptions.TLSConfig != nil && !dialOptions.TLSConfig.InsecureSkipVerify,
		"dialOptions.TimeLimit":           dialOptions.TimeLimit,
		"dialOptions.URL":                 dialOptions.URL.String(),
	}).Debug("Connecting")

	conn := supererrors.ExceptFn(supererrors.W(auth.Bind(
		bindParameters,
		dialOptions,
	)))

	logger.WithFields(logrus.Fields{
		"searchArguments.Attributes": getFlags.searchArguments.Attributes.ToAttributeList(),
		"searchArguments.Filter":     getFlags.searchArguments.Filter.String(),
		"searchArguments.Path":       getFlags.searchArguments.Path,
	}).Debug("Querying")

	results := supererrors.ExceptFn(supererrors.W(client.Search(conn, getFlags.searchArguments)))

	switch getFlags.format {
	case "csv":
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

	case "yaml":
		maps := make([]map[string]any, len(results))
		for i, r := range results {
			maps[i] = make(map[string]any)
			for _, a := range attributes.Map(r).Keys() {
				maps[i][a.String()] = r[a]
			}
		}

		if len(maps) == 1 {
			supererrors.Except(yaml.NewEncoder(apputil.Stdout(), yaml.Indent(2)).Encode(maps[0]))
		} else {
			supererrors.Except(yaml.NewEncoder(apputil.Stdout(), yaml.Indent(2)).Encode(map[string]any{"Results": maps}))
		}

	default:
		maps := make([]map[string]any, len(results))
		for i, r := range results {
			maps[i] = make(map[string]any)
			for _, a := range attributes.Map(r).Keys() {
				maps[i][a.String()] = r[a]
			}
		}

		buffer := bytes.NewBuffer(nil)
		if len(maps) == 1 {
			supererrors.Except(yaml.NewEncoder(buffer, yaml.Indent(2)).Encode(maps[0]))
		} else {
			supererrors.Except(yaml.NewEncoder(buffer, yaml.Indent(2)).Encode(map[string]any{"Results": maps}))
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
