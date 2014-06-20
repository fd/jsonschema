package jsonschema

import (
	"encoding/json"
)

type exclusiveMaximumValidator struct {
	exclusive bool
}

func (v *exclusiveMaximumValidator) Setup(x interface{}, builder Builder) error {
	if y, ok := x.(bool); ok && y {
		v.exclusive = true
	}

	return nil
}

func (v *exclusiveMaximumValidator) Validate(x interface{}, ctx *Context) {
	_, ok := x.(json.Number)
	if !ok {
		return
	}

	ctx.ExclusiveMaximum = v.exclusive
}
