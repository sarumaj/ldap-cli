package util

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var englishToTitleNoLower = cases.Title(language.English, cases.NoLower)

// Convert English words to title case (prevent lowercasing)
func ToTitleNoLower(in string) string { return englishToTitleNoLower.String(in) }
