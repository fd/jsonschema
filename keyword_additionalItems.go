package jsonschema

import (
	"fmt"
)

type additionalItemsValidator struct {
	allowed bool
	item    *Schema
}

func (v *additionalItemsValidator) Setup(x interface{}, e *Env) error {
	switch y := x.(type) {

	case map[string]interface{}:
		s, err := e.BuildSchema(y)
		if err != nil {
			return err
		}
		v.item = s
		v.allowed = true
		return nil

	case bool:
		v.allowed = y
		return nil

	default:
		return fmt.Errorf("invalid 'additionalItems' definition: %#v", x)

	}
}

func (v *additionalItemsValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return
	}

	if !v.allowed {
		if ctx.NextItem == 0 { // 'items' was not defined
			return
		}

		for i, l := ctx.NextItem, len(y); i < l; i++ {
			ctx.Report(&ErrInvalidItem{i, fmt.Errorf("additional item is not allowed")})
			ctx.NextItem = i + 1
		}
		return
	}

	if v.allowed && v.item != nil {
		for i, l := ctx.NextItem, len(y); i < l; i++ {
			err := v.item.Validate(y[i])
			if err != nil {
				ctx.Report(&ErrInvalidItem{i, err})
			}
			ctx.NextItem = i + 1
		}
		return
	}
}
