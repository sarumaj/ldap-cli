package util

import (
	"fmt"

	filter "github.com/sarumaj/ldap-cli/pkg/lib/definitions/filter"
)

func HelpAliases(msg *string) {
	if msg == nil {
		return
	}

	aliases := filter.ListAliases()

	*msg += "Filter option supports following alias expressions:\n"
	var longest int
	for _, alias := range aliases {
		if l := len(alias.String()); l > longest {
			longest = l
		}
	}

	for _, alias := range aliases {
		*msg += fmt.Sprintf(
			fmt.Sprintf(" - %%-%ds: %%s\n", longest),
			alias.String(), alias.Substitution(alias.Parameters),
		)
	}
	*msg += "Parameters beginning with a question mark are optional.\n"
}
