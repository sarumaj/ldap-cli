package util

import (
	"fmt"
	"net"
	"strings"
)

// Resolve Internet Protocol address to domain name
func LookupAddress(address string) string {
	ip, port, found := strings.Cut(address, ":")

	var names []string
	var err error
	if found {
		names, err = net.LookupAddr(ip)
	} else {
		names, err = net.LookupAddr(address)
	}

	if err == nil && len(names) > 0 {
		return fmt.Sprintf("%s:%s", strings.Trim(names[0], "."), port)
	}

	return address
}
