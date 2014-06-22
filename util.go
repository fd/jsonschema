package jsonschema

import (
	"encoding/json"
	"math"
	"net/url"
	"reflect"
	"strings"
)

func isEqual(a, b interface{}) (bool, error) {

	// handle numbers mathematically
	if n, ok := a.(json.Number); ok {
		if m, ok := b.(json.Number); ok {
			f, err := n.Float64()
			if err != nil {
				return false, err
			}

			g, err := m.Float64()
			if err != nil {
				return false, err
			}

			d := math.Abs(f - g)
			return d < 0.000000001, nil
		}
	}

	// other values
	return reflect.DeepEqual(a, b), nil

}

func isRef(x interface{}) (string, bool) {
	m, ok := x.(map[string]interface{})
	if !ok {
		return "", false
	}

	refi, found := m["$ref"]
	if !found {
		return "", false
	}

	ref, ok := refi.(string)
	if !ok {
		return "", false
	}

	if strings.IndexByte(ref, '#') < 0 {
		ref += "#"
	}

	return ref, true
}

func escapeJSONPointer(s string) string {
	s = strings.Replace(s, "~", "~0", -1)
	s = strings.Replace(s, "/", "~1", -1)
	return s
}

func normalizeRef(r string) string {
	if strings.IndexByte(r, '#') < 0 {
		r += "#"
	}
	return r
}

func resolveRef(base, ref *url.URL) *url.URL {
	dst := base.ResolveReference(ref)
	dst.Fragment = ref.Fragment
	return dst
}
