package jsonschema

import (
	"fmt"
	"reflect"
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

func (v *additionalItemsValidator) Validate(x reflect.Value, ctx *Context) {
	if !isArray(x) {
		return
	}

	if !v.allowed {
		if ctx.NextItem == 0 { // 'items' was not defined
			return
		}

		for i, l := ctx.NextItem, x.Len(); i < l; i++ {
			ctx.Report(&ItemValidationError{i, fmt.Errorf("additional item is not allowed")})
			ctx.NextItem = i + 1
		}
		return
	}

	if v.allowed && v.item != nil {
		for i, l := ctx.NextItem, x.Len(); i < l; i++ {
			err := v.item.ValidateValue(x.Index(i))
			if err != nil {
				ctx.Report(&ItemValidationError{i, err})
			}
			ctx.NextItem = i + 1
		}
		return
	}
}
