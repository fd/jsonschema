package jsonschema

import (
	"encoding/json"
	"fmt"
)

type ErrTooLarge struct {
	max       float64
	exclusive bool
	was       interface{}
}

func (e *ErrTooLarge) Error() string {
	if e.exclusive {
		return fmt.Sprintf("expected %#v to be smaller than %v", e.was, e.max)
	} else {
		return fmt.Sprintf("expected %#v to be smaller than or equal to %v", e.was, e.max)
	}
}

type maximumValidator struct {
	max float64
}

func (v *maximumValidator) Setup(x interface{}, builder Builder) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'maximum' definition: %#v", x)
	}

	f, err := y.Float64()
	if err != nil {
		return fmt.Errorf("invalid 'maximum' definition: %#v (%s)", x, err)
	}

	v.max = f
	return nil
}

func (v *maximumValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(json.Number)
	if !ok {
		return
	}

	f, err := y.Float64()
	if err != nil {
		ctx.Report(err)
		return
	}

	if ctx.ExclusiveMaximum {
		ok = f < v.max
	} else {
		ok = f <= v.max
	}

	if !ok {
		ctx.Report(&ErrTooLarge{v.max, ctx.ExclusiveMaximum, x})
	}
}
