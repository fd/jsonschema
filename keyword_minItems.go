package jsonschema

import (
	"encoding/json"
	"fmt"
)

type minItemsValidator struct {
	min int
}

func (v *minItemsValidator) Setup(x interface{}, e *Env) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'minItems' definition: %#v", x)
	}

	i, err := y.Int64()
	if err != nil {
		return fmt.Errorf("invalid 'minItems' definition: %#v (%s)", x, err)
	}

	v.min = int(i)
	return nil
}

func (v *minItemsValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return
	}

	if len(y) < v.min {
		ctx.Report(&ErrTooShort{v.min, x})
	}
}
