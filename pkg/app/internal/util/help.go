package util

import (
	"fmt"
	"strings"

	filter "github.com/sarumaj/ldap-cli/v2/pkg/lib/definitions/filter"
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
		subst := alias.Substitution(alias.Parameters)

		// mitigate improper formatting of the substitution
		for key, value := range map[string]string{
			"=<?operator>": "<?operator>",
			"<Attribute>":  "<attribute>",
		} {
			subst = strings.Replace(subst, key, value, 1)
		}

		*msg += fmt.Sprintf(
			fmt.Sprintf(" - %%-%ds: %%s\n", longest),
			alias.String(), subst,
		)
	}

	*msg += "Parameters beginning with a question mark are optional.\n"
}
