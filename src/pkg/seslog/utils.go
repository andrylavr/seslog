package seslog

import (
	"bytes"
	"encoding/binary"
	"net"
	"strings"
)

func IP2UInt32(ip net.IP) uint32 {
	var long uint32
	if err := binary.Read(bytes.NewBuffer(ip.To4()), binary.BigEndian, &long); err != nil {
		return 0
	}
	return long
}

func splitRequestURI(request_uri string) (string, string) {
	const c = "?"
	i := strings.Index(request_uri, c)
	if i < 0 {
		return request_uri, ""
	}
	return request_uri[:i], request_uri[i+len(c):]
}
