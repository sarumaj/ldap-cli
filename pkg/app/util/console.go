package util

import (
	"fmt"
	"os"
	"strings"

	term "golang.org/x/term"
)

// Is256ColorSupported checks if current terminal supports ANSI 256-bit color codes.
// Env variables TERM and COLORTERM are considered
func Is256ColorSupported() bool {
	return IsTrueColorSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

// IsColorEnabled checks if current terminal supports ANSI color codes.
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

// IsTerminal checks if given file descriptor is a terminal or just a pipe
func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// IsTrueColorSupported checks if current terminal supports ANSI 24-bit color codes.
// Env variables TERM and COLORTERM are evaluated
func IsTrueColorSupported() bool {
	spec := os.Getenv("TERM")
	colorSpec := os.Getenv("COLORTERM")

	return strings.Contains(spec, "24bit") ||
		strings.Contains(spec, "truecolor") ||
		strings.Contains(colorSpec, "24bit") ||
		strings.Contains(colorSpec, "truecolor")
}

// PrintColors prints formatted string with colors if supported
func PrintColors(fn func(string, ...any) string, format string, a ...any) string {
	if IsColorEnabled() {
		return fn(format, a...)
	}

	return fmt.Sprintf(format, a...)
}

// Stderr returns standard error file descriptor
func Stderr() *os.File { return os.Stderr }

// Stdin returns standard input file descriptor
func Stdin() *os.File { return os.Stdin }

// Stdout returns standard output file descriptor
func Stdout() *os.File { return os.Stdout }

// TerminalSize returns terminal size in columns and rows
func TerminalSize(f *os.File) (int, int, error) {
	return term.GetSize(int(f.Fd()))
}
