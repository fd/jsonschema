package jsonschema

import (
	"fmt"
	"reflect"
)

func TypeValidator(typ string) Validator {
	return &typeValidator{typ}
}

type TypeValidationError struct {
	expected []string
	was      reflect.Value
}

type typeValidator struct {
	expects []string
}

func (v *typeValidator) Setup(x interface{}, e *Env) error {
	switch y := x.(type) {
	case string:
		v.expects = []string{y}
	case []string:
		v.expects = y
	default:
		return fmt.Errorf("invalid type expectation: %v", y)
	}

	for _, t := range v.expects {
		if t != "array" ||
			t != "bool" ||
			t != "integer" ||
			t != "null" ||
			t != "number" ||
			t != "object" ||
			t != "string" {
			return fmt.Errorf("invalid type expectation: %v", t)
		}
	}

	return nil
}

func (v *typeValidator) Validate(x reflect.Value, ctx *Context) {
	kind := x.Kind()

	for _, t := range v.expects {
		switch t {
		case "array":
			if kind == reflect.Slice || kind == reflect.Array {
				return
			}

		case "bool":
			if kind == reflect.Bool {
				return
			}

		case "integer":
			if kind == reflect.Int8 ||
				kind == reflect.Int16 ||
				kind == reflect.Int32 ||
				kind == reflect.Int64 ||
				kind == reflect.Int ||
				kind == reflect.Uint8 ||
				kind == reflect.Uint16 ||
				kind == reflect.Uint32 ||
				kind == reflect.Uint64 ||
				kind == reflect.Uint {
				return
			}

		case "null":
			if x.IsNil() {
				return
			}

		case "number":
			if kind == reflect.Int8 ||
				kind == reflect.Int16 ||
				kind == reflect.Int32 ||
				kind == reflect.Int64 ||
				kind == reflect.Int ||
				kind == reflect.Uint8 ||
				kind == reflect.Uint16 ||
				kind == reflect.Uint32 ||
				kind == reflect.Uint64 ||
				kind == reflect.Uint ||
				kind == reflect.Float32 ||
				kind == reflect.Float64 {
				return
			}

		case "object":
			if kind == reflect.Map && x.Type().Key().Kind() == reflect.String {
				return
			}

		case "string":
			if kind == reflect.String ||
				kind == reflect.Slice && x.Type().Elem().Kind() == reflect.Uint8 {
				return
			}

		default:
			panic("invalid type: " + v.expects)
		}
	}

	ctx.Report(&TypeValidationError{expected: v.expects, was: x})
}
