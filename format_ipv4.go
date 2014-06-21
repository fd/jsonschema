package jsonschema

import (
	"net"
)

type ipv4Format struct{}

func (*ipv4Format) IsValid(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return true
	}

	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil
}
