package ip

import (
	"strings"
)

func GetIP(addr string) string {
	if strings.Count(addr, ":") > 1 {
		if sq := strings.LastIndexByte(addr, ']'); sq > 1 {
			return addr[1:sq]
		}

		return addr
	}

	host, _, found := strings.CutLast(addr, ":")
	if !found {
		return addr
	}

	return host
}
