package util

import (
	"context"
	"net"
	"strings"
)

// Resolve Internet Protocol address to domain name
func LookupAddress(address string) string {
	ip, port, found := strings.Cut(address, ":")

	var names []string
	var err error
	if ctx := context.Background(); found {
		names, err = net.DefaultResolver.LookupAddr(ctx, strings.Trim(ip, "[]"))
	} else {
		names, err = net.DefaultResolver.LookupAddr(ctx, address)
	}

	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(strings.Trim(names[0], ".")+":"+port, ":")
	}

	return address
}
