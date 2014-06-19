package jsonschema

import (
	"reflect"
)

type exclusiveMinimumValidator struct {
	exclusive bool
}

func (v *exclusiveMinimumValidator) Setup(x interface{}, e *Env) error {
	if y, ok := x.(bool); ok && y {
		v.exclusive = true
	}

	return nil
}

func (v *exclusiveMinimumValidator) Validate(x reflect.Value, ctx *Context) {
	if !isInteger(x) && !isFloat(x) {
		return
	}

	ctx.ExclusiveMinimum = v.exclusive
}
