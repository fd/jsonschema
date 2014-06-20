package jsonschema

import (
	"fmt"
)

type ErrInvalidItem struct {
	Index int
	Err   error
}

func (e *ErrInvalidItem) Error() string {
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

func (v *itemsValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return
	}

	if v.item != nil {
		for i, l := 0, len(y); i < l; i++ {
			err := v.item.Validate(y[i])
			if err != nil {
				ctx.Report(&ErrInvalidItem{i, err})
			}
			ctx.NextItem = i + 1
		}
		return
	}

	if len(v.items) > 0 {
		for i, la, lb := 0, len(y), len(v.items); i < la && i < lb; i++ {
			err := v.items[i].Validate(y[i])
			if err != nil {
				ctx.Report(&ErrInvalidItem{i, err})
			}
			ctx.NextItem = i + 1
		}
		return
	}
}
