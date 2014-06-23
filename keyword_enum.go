package jsonschema

import (
	"fmt"
)

type enumValidator struct {
	enum []interface{}
}

func (v *enumValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("enum"); found {
		y, ok := x.([]interface{})
		if !ok || y == nil || len(y) == 0 {
			return fmt.Errorf("invalid 'enum' definition: %#v", x)
		}

		v.enum = y
	}
	return nil
}

func (v *enumValidator) Validate(x interface{}, ctx *Context) {
	for _, y := range v.enum {
		equal, err := isEqual(x, y)
		if err != nil {
			ctx.Report(err)
		}
		if equal {
			return
		}
	}

	ctx.Report(&ErrInvalidEnum{v.enum, x})
}
