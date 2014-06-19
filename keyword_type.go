package jsonschema

import (
	"fmt"
	"reflect"
)

type TypeValidationError struct {
	expected []PrimitiveType
	was      reflect.Value
}

func (e *TypeValidationError) Error() string {
	return fmt.Sprintf("expected type to be in %#v but was %#v", e.expected, e.was.Interface())
}

type typeValidator struct {
	expects []PrimitiveType
}

func (v *typeValidator) Setup(x interface{}, e *Env) error {
	switch y := x.(type) {
	case string:
		v.expects = []PrimitiveType{PrimitiveType(y)}

	case []string:
		var z = make([]PrimitiveType, len(y))
		for i, a := range y {
			z[i] = PrimitiveType(a)
		}
		v.expects = z

	case []interface{}:
		var z = make([]PrimitiveType, len(y))
		for i, a := range y {
			if b, ok := a.(string); ok {
				z[i] = PrimitiveType(b)
			} else {
				return fmt.Errorf("invalid type expectation: %#v", y)
			}
		}
		v.expects = z

	default:
		return fmt.Errorf("invalid type expectation: %#v", y)
	}

	for _, t := range v.expects {
		if !t.Valid() {
			return fmt.Errorf("invalid type expectation: %#v", t)
		}
	}

	return nil
}

func (v *typeValidator) Validate(x reflect.Value, ctx *Context) {
	for x.Kind() == reflect.Interface || x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	kind := x.Kind()

	for _, t := range v.expects {
		switch t {
		case ArrayType:
			if isArray(x) {
				ctx.Type = ArrayType
				return
			}

		case BooleanType:
			if kind == reflect.Bool {
				ctx.Type = BooleanType
				return
			}

		case IntegerType:
			if isInteger(x) {
				ctx.Type = IntegerType
				return
			}
			if isIntegerFloat(x) {
				ctx.Type = IntegerType
				return
			}

		case NullType:
			if (kind == reflect.Ptr ||
				kind == reflect.Interface ||
				kind == reflect.Map ||
				kind == reflect.Slice) && x.IsNil() {
				ctx.Type = NullType
				return
			}

		case NumberType:
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
				ctx.Type = NumberType
				return
			}

		case ObjectType:
			if kind == reflect.Map && x.Type().Key().Kind() == reflect.String ||
				kind == reflect.Struct {
				ctx.Type = ObjectType
				return
			}

		case StringType:
			if isString(x) {
				ctx.Type = StringType
				return
			}

		default:
			panic("invalid type: " + t)
		}
	}

	ctx.Report(&TypeValidationError{expected: v.expects, was: x})
}
