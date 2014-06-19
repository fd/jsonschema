package jsonschema

import (
	"fmt"
	"reflect"
)

type ItemValidationError struct {
	Index int
	Err   error
}

func (e *ItemValidationError) Error() string {
	return fmt.Sprintf("Invalid item at %v: %s", e.Index, e.Err)
}

type itemsValidator struct {
	item  *Schema
	items []*Schema
}

func (v *itemsValidator) Setup(x interface{}, e *Env) error {
	switch y := x.(type) {

	case map[string]interface{}:
		s, err := e.BuildSchema(y)
		if err != nil {
			return err
		}
		v.item = s
		return nil

	case []interface{}:
		l := make([]*Schema, len(y))
		for i, a := range y {
			b, ok := a.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid 'items' definition: %#v", x)
			}
			s, err := e.BuildSchema(b)
			if err != nil {
				return err
			}
			l[i] = s
		}
		v.items = l
		return nil

	default:
		return fmt.Errorf("invalid 'items' definition: %#v", x)

	}
}

func (v *itemsValidator) Validate(x reflect.Value, ctx *Context) {
	if !isArray(x) {
		return
	}

	if v.item != nil {
		for i, l := 0, x.Len(); i < l; i++ {
			err := v.item.ValidateValue(x.Index(i))
			if err != nil {
				ctx.Report(&ItemValidationError{i, err})
			}
			ctx.NextItem = i + 1
		}
		return
	}

	if len(v.items) > 0 {
		for i, la, lb := 0, x.Len(), len(v.items); i < la && i < lb; i++ {
			err := v.items[i].ValidateValue(x.Index(i))
			if err != nil {
				ctx.Report(&ItemValidationError{i, err})
			}
			ctx.NextItem = i + 1
		}
		return
	}
}
