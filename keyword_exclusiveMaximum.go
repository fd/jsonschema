package jsonschema

import (
	"reflect"
)

type exclusiveMaximumValidator struct {
	exclusive bool
}

func (v *exclusiveMaximumValidator) Setup(x interface{}, e *Env) error {
	if y, ok := x.(bool); ok && y {
		v.exclusive = true
	}

	return nil
}

func (v *exclusiveMaximumValidator) Validate(x reflect.Value, ctx *Context) {
	if !isInteger(x) && !isFloat(x) {
		return
	}

	ctx.ExclusiveMaximum = v.exclusive
}
