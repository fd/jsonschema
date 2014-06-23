package jsonschema

import (
	"fmt"
)

type maximumValidator struct {
	max       float64
	exclusive bool
}

func (v *maximumValidator) Setup(builder Builder) error {
	if x, ok := builder.GetKeyword("exclusiveMaximum"); ok {
		y, ok := x.(bool)

		if !ok {
			return fmt.Errorf("invalid 'exclusiveMaximum' definition: %#v", x)
		}

		v.exclusive = y
	}

	if x, found := builder.GetKeyword("maximum"); found {
		f, ok, err := toFloat(x)
		if !ok {
			return fmt.Errorf("invalid 'maximum' definition: %#v", x)
		}
		if err != nil {
			return fmt.Errorf("invalid 'maximum' definition: %#v (%s)", x, err)
		}

		v.max = f
	}

	return nil
}

func (v *maximumValidator) Validate(x interface{}, ctx *Context) {
	f, ok, err := toFloat(x)
	if !ok {
		return
	}
	if err != nil {
		ctx.Report(err)
		return
	}

	if v.exclusive {
		ok = f < v.max
	} else {
		ok = f <= v.max
	}

	if !ok {
		ctx.Report(&ErrTooLarge{v.max, v.exclusive, x})
	}
}
