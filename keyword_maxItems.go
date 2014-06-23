package jsonschema

import (
	"fmt"
)

type maxItemsValidator struct {
	max int
}

func (v *maxItemsValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("maxItems"); found {
		i, ok := x.(int64)
		if !ok {
			return fmt.Errorf("invalid 'maxItems' definition: %#v", x)
		}

		v.max = int(i)
	}
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
