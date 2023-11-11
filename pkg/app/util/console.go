package util

import (
	"fmt"
	"os"
	"strings"

	term "golang.org/x/term"
)

// Check if current terminal supports ANSI 256-bit color codes.
// Env variables TERM and COLORTERM are considered
func Is256ColorSupported() bool {
	return IsTrueColorSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

// Check if current terminal has colors enabled.
// Env variables CLICOLOR_FORCE, NO_COLOR and CLICOLOR are being considered
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

// Check if given terminal descriptor is a terminal or just a pipe
func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// Check fif terminal supports true color mode.
// Env variables TERM and COLORTERM are evaluated
func IsTrueColorSupported() bool {
	spec := os.Getenv("TERM")
	colorSpec := os.Getenv("COLORTERM")

	return strings.Contains(spec, "24bit") ||
		strings.Contains(spec, "truecolor") ||
		strings.Contains(colorSpec, "24bit") ||
		strings.Contains(colorSpec, "truecolor")
}

// Print colorized text only if terminal supports ANSI color codes
func PrintColors(fn func(string, ...any) string, format string, a ...any) string {
	if IsColorEnabled() {
		return fn(format, a...)
	}

	return fmt.Sprintf(format, a...)
}

func Stderr() *os.File { return os.Stderr } // Standard destination for errors
func Stdin() *os.File  { return os.Stdin }  // Standard input
func Stdout() *os.File { return os.Stdout } // Standard output

// Get size of current terminal window
func TerminalSize(f *os.File) (int, int, error) {
	return term.GetSize(int(f.Fd()))
}
