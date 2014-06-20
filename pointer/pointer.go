// JSON Pointer implemenation
//
// See: http://tools.ietf.org/html/draft-ietf-appsawg-json-pointer-04

package pointer

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidPointer = errors.New("invalid json-pointer")

type Pointer struct {
	components []string
}

func Parse(p string) (Pointer, error) {
	if len(p) == 0 {
		return Pointer{[]string{}}, nil
	}

	if p[0] != '/' {
		return Pointer{}, ErrInvalidPointer
	}

	l := strings.Split(p[1:], "/")

	for i, c := range l {
		// validate escapes
		component_length := len(c)
		for j, idx := 0, strings.IndexByte(c[j:], '~'); idx >= 0; j += idx + 1 {
			tilda_index := j + idx
			code_index := tilda_index + 1
			if tilda_index >= component_length {
				return Pointer{}, ErrInvalidPointer
			}
			code := c[code_index]
			if code != '0' && code != '1' {
				return Pointer{}, ErrInvalidPointer
			}
		}

		c = strings.Replace(strings.Replace(c, "~1", "/", -1), "~0", "~", -1)

		l[i] = c
	}

	return Pointer{l}, nil
}

func (p Pointer) Find(v interface{}) (interface{}, bool) {
	for _, c := range p.components {
		switch x := v.(type) {

		case []interface{}:
			idx, err := strconv.ParseUint(c, 10, 64)
			if err != nil {
				return nil, false
			}
			if idx >= len(x) {
				return nil, false
			}
			v = x[idx]

		case map[string]interface{}:
			m, found := x[c]
			if !found {
				return nil, false
			}
			v = m

		default:
			return nil, false

		}
	}

	return v, true
}
