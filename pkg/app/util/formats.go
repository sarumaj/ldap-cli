package util

import (
	"fmt"
	"slices"
)

const (
	CSV     = "csv"
	DEFAULT = "default"
	LDIF    = "ldif"
	YAML    = "yaml"
)

var supportedFormats = []string{CSV, DEFAULT, LDIF, YAML}

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
