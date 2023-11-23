package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"slices"

	color "github.com/fatih/color"
	ldif "github.com/go-ldap/ldif"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/sarumaj/ldap-cli/pkg/lib/definitions/attributes"
)

const (
	CSV     = "csv"
	DEFAULT = "default"
	LDIF    = "ldif"
	YAML    = "yaml"
)

var supportedFormats = []string{CSV, DEFAULT, LDIF, YAML}

func FlushToStdOut(results attributes.Maps, requests *ldif.LDIF, format string) error {
	switch format {

	case CSV:
		lines := make([][]string, len(results)+1)
		for i, m := range results {
			for _, a := range attributes.Map(m).Keys() {
				if i == 0 {
					lines[0] = append(lines[0], a.String())
				}

				lines[i+1] = append(lines[i+1], fmt.Sprint(m[a]))
			}
		}
		return csv.NewWriter(Stdout()).WriteAll(lines)

	case LDIF:
		data, err := ldif.Marshal(requests)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(Stdout(), data)
		return err

	case YAML:
		if len(results) == 1 {
			return yaml.NewEncoder(Stdout(), yaml.Indent(2)).Encode(results[0])
		} else {
			return yaml.NewEncoder(Stdout(), yaml.Indent(2)).Encode(map[string]any{"Results": results})
		}

	default:
		buffer := bytes.NewBuffer(nil)

		var err error
		if len(results) == 1 {
			err = yaml.NewEncoder(buffer, yaml.Indent(2)).Encode(results[0])
		} else {
			err = yaml.NewEncoder(buffer, yaml.Indent(2)).Encode(map[string]any{"Results": results})
		}
		if err != nil {
			return err
		}

		tokens := lexer.Tokenize(buffer.String())
		buffer.Reset()

		_, err = fmt.Fprintln(Stdout(), (&printer.Printer{
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
		}).PrintTokens(tokens))

		return err
	}
}

func ListSupportedFormats(quote bool) (list []string) {
	for _, f := range supportedFormats {
		if quote {
			list = append(list, fmt.Sprintf("%q", f))
		} else {
			list = append(list, f)
		}
	}

	slices.Sort(list)
	return
}
