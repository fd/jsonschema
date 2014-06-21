package jsonschema

import (
	"encoding/json"
	"fmt"
)

type ErrInvalidType struct {
	expected []PrimitiveType
	was      interface{}
}

func (e *ErrInvalidType) Error() string {
	return fmt.Sprintf("expected type to be in %#v but was %#v", e.expected, e.was)
}

type typeValidator struct {
	expects []PrimitiveType
}

func (v *typeValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("type"); found {
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
					return fmt.Errorf("invalid type expectation: %#v", x)
				}
			}
			v.expects = z

		default:
			return fmt.Errorf("invalid type expectation: %#v", x)
		}

		for _, t := range v.expects {
			if !t.Valid() {
				return fmt.Errorf("invalid type expectation: %#v", t)
			}
		}
	}
	return nil
}

func (v *typeValidator) Validate(x interface{}, ctx *Context) {
	for _, t := range v.expects {
		switch t {
		case ArrayType:
			if _, ok := x.([]interface{}); ok && x != nil {
				ctx.Type = ArrayType
				return
			}

		case BooleanType:
			if _, ok := x.(bool); ok {
				ctx.Type = BooleanType
				return
			}

		case IntegerType:
			if y, ok := x.(json.Number); ok {
				_, err := y.Int64()
				if err == nil {
					ctx.Type = IntegerType
					return
				}
			}

		case NullType:
			if x == nil {
				ctx.Type = NullType
				return
			}

		case NumberType:
			if y, ok := x.(json.Number); ok {
				_, err := y.Float64()
				if err == nil {
					ctx.Type = NumberType
					return
				}
			}

		case ObjectType:
			if _, ok := x.(map[string]interface{}); ok && x != nil {
				ctx.Type = ObjectType
				return
			}

		case StringType:
			if _, ok := x.(string); ok {
				ctx.Type = StringType
				return
			}

		default:
			panic("invalid type: " + t)
		}
	}

	ctx.Report(&ErrInvalidType{expected: v.expects, was: x})
}
