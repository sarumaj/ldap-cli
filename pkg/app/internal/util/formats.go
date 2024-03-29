package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	color "github.com/fatih/color"
	ldif "github.com/go-ldap/ldif"
	yaml "github.com/goccy/go-yaml"
	lexer "github.com/goccy/go-yaml/lexer"
	printer "github.com/goccy/go-yaml/printer"
	attributes "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/attributes"
)

const (
	CSV     = "csv"     // CSV format, "," being used as delimiter
	DEFAULT = "default" // Colorful YAML output
	LDIF    = "ldif"    // LDAP Data Interchange Format
	YAML    = "yaml"    // YAML
)

// Config used to colorize YAML output (default)
var defaultPrinterConfig = &printer.Printer{
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
}

// List of supported formats
var supportedFormats = []string{CSV, DEFAULT, LDIF, YAML}

// Encode results into given format.
// Writer is supposed to be either stdout or a file.
// Per default, colorful YAML format takes precedences and is being emitted to stdout
func Flush(results attributes.Maps, requests *ldif.LDIF, format string, out io.Writer) error {
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
		return csv.NewWriter(out).WriteAll(lines)

	case LDIF:
		data, err := ldif.Marshal(requests)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, data)
		return err

	case YAML:
		if len(results) == 1 {
			return yaml.NewEncoder(out, yaml.Indent(2)).Encode(results[0])
		} else {
			return yaml.NewEncoder(out, yaml.Indent(2)).Encode(map[string]any{"Results": results})
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

		if !IsColorEnabled() {
			_, err = fmt.Fprintln(Stdout(), buffer.String())
			return err
		}

		tokens := lexer.Tokenize(buffer.String())
		buffer.Reset()

		_, err = fmt.Fprintln(Stdout(), defaultPrinterConfig.PrintTokens(tokens))

		return err
	}
}

// ListSupportedFormats returns a list of supported formats
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

// SniffFormat returns the format of a file from its extension
func SniffFormat(filename, format string) string {
	switch format := strings.TrimPrefix(filepath.Ext(filename), "."); format {

	case CSV, LDIF, YAML:
		return format

	}

	return format
}
