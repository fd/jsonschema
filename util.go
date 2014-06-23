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
	if f, ok, err := toFloat(a); ok {
		if err != nil {
			return false, err
		}
		if g, ok, err := toFloat(b); ok {
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

func toFloat(x interface{}) (float64, bool, error) {
	switch y := x.(type) {

	case json.Number:
		f, err := y.Float64()
		if err != nil {
			return 0, true, err
		}
		return f, true, nil

	case int64:
		return float64(y), true, nil

	case float64:
		return y, true, nil

	default:
		return 0, false, nil

	}
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

func rootRef(ref string) string {
	ref = normalizeRef(ref)
	idx := strings.IndexByte(ref, '#')
	return ref[:idx+1]
}

func refURL(ref string) string {
	ref = normalizeRef(ref)
	idx := strings.IndexByte(ref, '#')
	return ref[:idx]
}

func refFragment(ref string) string {
	ref = normalizeRef(ref)
	idx := strings.IndexByte(ref, '#')
	return ref[idx+1:]
}
