package util

import (
	"fmt"
	"net"
	"strings"
)

func Resolve(remoteHost string) string {
	ip, port, found := strings.Cut(remoteHost, ":")

	var names []string
	var err error
	if found {
		names, err = net.LookupAddr(ip)
	} else {
		names, err = net.LookupAddr(remoteHost)
	}

	if err == nil && len(names) > 0 {
		return fmt.Sprintf("%s:%s", strings.Trim(names[0], "."), port)
	}

	return remoteHost
}
