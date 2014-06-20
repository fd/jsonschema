package jsonschema

import (
	"encoding/json"
	"fmt"
)

type minPropertiesValidator struct {
	min int
}

func (v *minPropertiesValidator) Setup(x interface{}, e *Env) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'minProperties' definition: %#v", x)
	}

	i, err := y.Int64()
	if err != nil {
		return fmt.Errorf("invalid 'minProperties' definition: %#v (%s)", x, err)
	}

	v.min = int(i)
	return nil
}

func (v *minPropertiesValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return
	}

	l := len(y)

	if l < v.min {
		ctx.Report(&ErrTooLong{v.min, x})
	}
}
