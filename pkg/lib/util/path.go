package util

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
)

// GetExecutablePath returns the path to the current executable
func GetExecutablePath() string {
	executablePath, err := os.Executable()
	if err != nil {
		return "unknown"
	}

	evaluatedPath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		return executablePath
	}

	return evaluatedPath
}

// RebuildStringSliceFlag rebuilds a string slice flag from a string by using custom CSV reader
func RebuildStringSliceFlag(flags []string, delimiter rune) ([]string, error) {
	in := strings.Join(flags, string(delimiter))
	reader := csv.NewReader(strings.NewReader(in))
	reader.Comma = delimiter
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	rebuilt, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var list []string
	for _, fragment := range rebuilt {
		if len(fragment) > 0 {
			list = append(list, fragment)
		}
	}

	return list, nil
}
