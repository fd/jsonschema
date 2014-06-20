package jsonschema

import (
	"encoding/json"
	"fmt"
)

type maxItemsValidator struct {
	max int
}

func (v *maxItemsValidator) Setup(x interface{}, builder Builder) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'maxItems' definition: %#v", x)
	}

	i, err := y.Int64()
	if err != nil {
		return fmt.Errorf("invalid 'maxItems' definition: %#v (%s)", x, err)
	}

	v.max = int(i)
	return nil
}

func (v *maxItemsValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return
	}

	if len(y) > v.max {
		ctx.Report(&ErrTooLong{v.max, x})
	}
}
