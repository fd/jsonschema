package jsonschema

import (
	"encoding/json"
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

func (v *exclusiveMinimumValidator) Validate(x interface{}, ctx *Context) {
	_, ok := x.(json.Number)
	if !ok {
		return
	}

	ctx.ExclusiveMinimum = v.exclusive
}
