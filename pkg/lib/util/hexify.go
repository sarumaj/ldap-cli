package util

import (
	"fmt"
	"strconv"
	"strings"
)

// Hexify converts a string to a hex representation
func Hexify(s string) string {
	var chars []string
	for _, c := range []byte(s) {
		chars = append(chars, fmt.Sprintf("\\x%02X", c))
	}

	return strings.Join(chars, "")
}

// Unhexify converts a hex representation to a string
func Unhexify(s string) string {
	chars := []byte{}
	for _, c := range strings.Split(strings.TrimLeft(s, "\\x"), "\\x") {
		i, _ := strconv.ParseInt(c, 16, 8)
		chars = append(chars, byte(i))
	}
	return string(chars)
}
