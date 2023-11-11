package commands

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	supererrors "github.com/sarumaj/go-super/errors"
	apputil "github.com/sarumaj/ldap-cli/pkg/app/util"
	auth "github.com/sarumaj/ldap-cli/pkg/lib/auth"
	client "github.com/sarumaj/ldap-cli/pkg/lib/client"
	attributes "github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
	"github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var filterString string
var searchArguments = &client.SearchArguments{}
var selectAttributes string
var format string

var searchCmd = func() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:     "search",
		Short:   "Search a directory object",
		Example: "ldap-cli search <object>",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			parent := cmd.Parent()
			parent.PersistentPreRun(parent, args)

			if searchArguments.Path == "" {
				var components []string
				for _, dc := range strings.Split(dialOptions.URL.Host, ".") {
					if dc == "" {
						continue
					}

					components = append(components, "DC="+dc)
				}
				searchArguments.Path = strings.Join(components, ",")
			}

			if len(selectAttributes) > 0 {
				reader := csv.NewReader(strings.NewReader(selectAttributes))
				reader.LazyQuotes = true
				reader.TrimLeadingSpace = true

				searchArguments.Attributes = attributes.LookupMany(supererrors.ExceptFn(supererrors.W(reader.Read()))...)
			}
		},
		Run: func(cmd *cobra.Command, _ []string) {
			supererrors.Except(cmd.Help())
		},
	}

	flags := searchCmd.PersistentFlags()
	flags.StringVar(&format, "format", "yaml", "Output format (supported: [\"csv\", \"yaml\"])")
	flags.StringVar(&searchArguments.Path, "path", "", "Specify the query path to search the directory objects in")
	flags.StringVar(&selectAttributes, "select", "", "Select specific object attributes")

	searchCmd.AddCommand(searchCustomCmd, searchUserCmd)

	return searchCmd
}()

func performSearch(cmd *cobra.Command, _ []string) {
	command := []string{cmd.Name()}
	for c := cmd.Parent(); c != nil; c.Parent() {
		command = append(command, c.Name())
	}

	logger := cmdLogger.WithField("command", command)
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
		"searchArguments.Attributes": searchArguments.Attributes.ToAttributeList(),
		"searchArguments.Filter":     searchArguments.Filter.String(),
		"searchArguments.Path":       searchArguments.Path,
	}).Debug("Querying")
	results := supererrors.ExceptFn(supererrors.W(client.Search(conn, *searchArguments)))

	switch format {
	case "csv":
		w := csv.NewWriter(apputil.Stdout())
		lines := make([][]string, len(results)+1)
		for i, m := range results {
			for _, a := range attributes.Map(m).Keys() {
				if i == 0 {
					lines[0] = append(lines[0], a.String())
				}
				lines[i+1] = append(lines[i+1], fmt.Sprint(m[a]))
			}
		}
		supererrors.Except(w.WriteAll(lines))

	default:
		maps := make([]map[string]any, len(results))
		for i, r := range results {
			maps[i] = make(map[string]any)
			for _, a := range attributes.Map(r).Keys() {
				maps[i][a.String()] = r[a]
			}
		}
		supererrors.Except(apputil.NewColoredYAMLEncoder(apputil.Stdout(), yaml.Indent(2)).Encode(maps))

	}

}
