package jsonschema

import (
	"strings"
)

// See:
//  http://tools.ietf.org/html/rfc1034#section-3.1
//  http://en.wikipedia.org/wiki/Domain_Name_System#Domain_name_syntax
type hostnameFormat struct{}

func (*hostnameFormat) IsValid(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return true
	}

	if len(s) > 253 || len(s) == 0 {
		return false
	}

	for len(s) > 0 {
		var (
			label string
			idx   = strings.IndexByte(s, '.')
		)

		if idx == 0 {
			return false
		} else if idx < 0 {
			label = s
			s = ""
		} else {
			label = s[:idx]
			s = s[idx+1:]
		}

		if len(label) > 63 || len(label) == 0 {
			return false
		}

		last_i := len(label) - 1
		for i, char := range label {
			if 'a' <= char && char <= 'z' {
				continue
			}
			if 'A' <= char && char <= 'Z' {
				continue
			}
			if '0' <= char && char <= '9' {
				continue
			}
			if i != 0 && i != last_i && char == '-' {
				continue
			}
			return false
		}
	}

	return true
}
