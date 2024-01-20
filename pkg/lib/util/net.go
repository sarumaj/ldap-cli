package util

import (
	"net"
	"strings"
)

// LookupAddress returns the hostname for the given address
func LookupAddress(address string) string {
	ip, port, found := strings.Cut(strings.Trim(address, "[]"), ":")

	var names []string
	var err error
	if found {
		names, err = net.LookupAddr(strings.Trim(ip, "[]"))
	} else {
		names, err = net.LookupAddr(address)
	}

	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(strings.Trim(names[0], ".")+":"+port, ":")
	}

	return address
}
