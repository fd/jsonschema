package jsonschema

import (
	"encoding/json"
	"fmt"
)

type ErrTooSmall struct {
	min       float64
	exclusive bool
	was       interface{}
}

func (e *ErrTooSmall) Error() string {
	if e.exclusive {
		return fmt.Sprintf("expected %#v to be larger than %v", e.was, e.min)
	} else {
		return fmt.Sprintf("expected %#v to be larger than or equal to %v", e.was, e.min)
	}
}

type minimumValidator struct {
	min float64
}

func (v *minimumValidator) Setup(x interface{}, builder Builder) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'minimum' definition: %#v", x)
	}

	f, err := y.Float64()
	if err != nil {
		return fmt.Errorf("invalid 'minimum' definition: %#v (%s)", x, err)
	}

	v.min = f
	return nil
}

func (v *minimumValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(json.Number)
	if !ok {
		return
	}

	f, err := y.Float64()
	if err != nil {
		ctx.Report(err)
		return
	}

	if ctx.ExclusiveMinimum {
		ok = f > v.min
	} else {
		ok = f >= v.min
	}

	if !ok {
		ctx.Report(&ErrTooSmall{v.min, ctx.ExclusiveMinimum, x})
	}
}
