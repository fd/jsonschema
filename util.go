package jsonschema

import (
	"math"
	"reflect"
)

func isArray(x reflect.Value) bool {
	kind := x.Kind()
	return (kind == reflect.Slice || kind == reflect.Array) && !x.IsNil()
}

func isInteger(x reflect.Value) bool {
	kind := x.Kind()
	return kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Int ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Uint
}

func isIntegerFloat(x reflect.Value) bool {
	if isFloat(x) {
		f := x.Float()
		return f == math.Trunc(f)
	}
	return false
}

func isFloat(x reflect.Value) bool {
	kind := x.Kind()
	return kind == reflect.Float32 || kind == reflect.Float64
}

func isString(x reflect.Value) bool {
	kind := x.Kind()
	return kind == reflect.String ||
		kind == reflect.Slice && x.Type().Elem().Kind() == reflect.Uint8
}
