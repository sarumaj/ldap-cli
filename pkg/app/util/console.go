package util

import (
	"fmt"
	"os"
	"strings"

	term "golang.org/x/term"
)

func CheckColors(fn func(string, ...any) string, format string, a ...any) string {
	if IsColorEnabled() {
		return fn(format, a...)
	}

	return fmt.Sprintf(format, a...)
}

func Is256ColorSupported() bool {
	return IsTrueColorSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

func IsColorEnabled() bool {
	switch {
	case
		// forced colored terminal output
		os.Getenv("CLICOLOR_FORCE") != "" && os.Getenv("CLICOLOR_FORCE") != "0",
		// not disabled and supported colored output
		os.Getenv("NO_COLOR") == "" && os.Getenv("CLICOLOR") != "0" && IsTerminal(Stdout()):

		return true
	}

	return false
}

func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func IsTrueColorSupported() bool {
	spec := os.Getenv("TERM")
	colorSpec := os.Getenv("COLORTERM")

	return strings.Contains(spec, "24bit") ||
		strings.Contains(spec, "truecolor") ||
		strings.Contains(colorSpec, "24bit") ||
		strings.Contains(colorSpec, "truecolor")
}

func Stderr() *os.File { return os.Stderr }
func Stdin() *os.File  { return os.Stdin }
func Stdout() *os.File { return os.Stdout }

func TerminalSize(f *os.File) (int, int, error) {
	return term.GetSize(int(f.Fd()))
}
