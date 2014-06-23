package jsonschema

import (
	"fmt"
)

type minimumValidator struct {
	min       float64
	exclusive bool
}

func (v *minimumValidator) Setup(builder Builder) error {
	if x, ok := builder.GetKeyword("exclusiveMinimum"); ok {
		y, ok := x.(bool)

		if !ok {
			return fmt.Errorf("invalid 'exclusiveMinimum' definition: %#v", x)
		}

		v.exclusive = y
	}

	if x, found := builder.GetKeyword("minimum"); found {
		f, ok, err := toFloat(x)
		if !ok {
			return fmt.Errorf("invalid 'minimum' definition: %#v", x)
		}
		if err != nil {
			return fmt.Errorf("invalid 'minimum' definition: %#v (%s)", x, err)
		}

		v.min = f
	}
	return nil
}

func (v *minimumValidator) Validate(x interface{}, ctx *Context) {
	f, ok, err := toFloat(x)
	if !ok {
		return
	}
	if err != nil {
		ctx.Report(err)
		return
	}

	if v.exclusive {
		ok = f > v.min
	} else {
		ok = f >= v.min
	}

	if !ok {
		ctx.Report(&ErrTooSmall{v.min, v.exclusive, x})
	}
}
