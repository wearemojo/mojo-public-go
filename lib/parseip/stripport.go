package parseip

import (
	"strings"
)

func StripPort(addr string) string {
	if strings.Count(addr, ":") > 1 {
		if sq := strings.LastIndexByte(addr, ']'); sq > 1 {
			return addr[1:sq]
		}

		return addr
	}

	i := strings.LastIndexByte(addr, ':')
	if i == -1 {
		return addr
	}

	return addr[:i]
}
