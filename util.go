package jsonschema

import (
	"encoding/json"
	"math"
	"reflect"
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
