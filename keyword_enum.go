package jsonschema

import (
	"fmt"
)

type ErrInvalidEnum struct {
	Expected []interface{}
	Value    interface{}
}

func (e *ErrInvalidEnum) Error() string {
	return fmt.Sprintf("%v must be in %v", e.Value, e.Expected)
}

type enumValidator struct {
	enum []interface{}
}

func (v *enumValidator) Setup(x interface{}, builder Builder) error {
	y, ok := x.([]interface{})
	if !ok || y == nil || len(y) == 0 {
		return fmt.Errorf("invalid 'enum' definition: %#v", x)
	}

	v.enum = y
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
