package jsonschema

import (
	"strings"
)

type emailFormat struct {
	hostname hostnameFormat
	ipv4     ipv4Format
	ipv6     ipv6Format
}

func (f *emailFormat) IsValid(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return true
	}

	if len(s) == 0 {
		return false
	}

	idx := 0
	if s[0] == '"' {
		offset := 1
		for {
			i := strings.IndexByte(s[offset:], '"')
			if i < 0 {
				return false
			}
			offset += i + 1
			if s[offset-2] != '\\' {
				break
			}
		}

		idx = strings.IndexByte(s[offset:], '@')
		if idx != 0 {
			return false
		}
		idx = offset + idx
	} else {
		idx = strings.IndexByte(s, '@')
		if idx <= 0 {
			return false
		}
	}

	node := s[idx+1:]
	if len(node) == 0 {
		return false
	}

	if strings.HasPrefix(node, "[IPv6:") && strings.HasSuffix(node, "]") {
		return f.ipv6.IsValid(node[6 : len(node)-1])
	}

	if strings.HasPrefix(node, "[") && strings.HasSuffix(node, "]") {
		return f.ipv6.IsValid(node[1 : len(node)-1])
	}

	return f.hostname.IsValid(node)
}
