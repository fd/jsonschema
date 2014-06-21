package jsonschema

import (
	"encoding/json"
	"fmt"
)

type maxPropertiesValidator struct {
	max int
}

func (v *maxPropertiesValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("maxProperties"); found {
		y, ok := x.(json.Number)
		if !ok {
			return fmt.Errorf("invalid 'maxProperties' definition: %#v", x)
		}

		i, err := y.Int64()
		if err != nil {
			return fmt.Errorf("invalid 'maxProperties' definition: %#v (%s)", x, err)
		}

		v.max = int(i)
	}
	return nil
}

func (v *maxPropertiesValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return
	}

	l := len(y)

	if l > v.max {
		ctx.Report(&ErrTooLong{v.max, x})
	}
}
