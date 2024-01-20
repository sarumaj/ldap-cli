package util

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// englishToTitle is a TitleCase converter for English language
var englishToTitleNoLower = cases.Title(language.English, cases.NoLower)

// ToTitleNoLower converts a string to title case without lowercasing
func ToTitleNoLower(in string) string { return englishToTitleNoLower.String(in) }
