package jsonschema

import (
	"fmt"
)

var additionalItemsDenied = &Schema{}

type itemsValidator struct {
	item           *Schema
	items          []*Schema
	additionalItem *Schema
}

func (v *itemsValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("items"); found {
		switch y := x.(type) {

		case map[string]interface{}:
			s, err := builder.Build("/items", y)
			if err != nil {
				return err
			}
			v.item = s

		case []interface{}:
			l := make([]*Schema, len(y))
			for i, a := range y {
				b, ok := a.(map[string]interface{})
				if !ok {
					return fmt.Errorf("invalid 'items' definition: %#v", y)
				}
				s, err := builder.Build(fmt.Sprintf("/items/%d", i), b)
				if err != nil {
					return err
				}
				l[i] = s
			}
			v.items = l

		default:
			return fmt.Errorf("invalid 'items' definition: %#v", y)

		}
	}

	if x, ok := builder.GetKeyword("additionalItems"); ok {
		switch y := x.(type) {

		case map[string]interface{}:
			s, err := builder.Build("/additionalItems", y)
			if err != nil {
				return err
			}
			v.additionalItem = s

		case bool:
			if !y {
				v.additionalItem = additionalItemsDenied
			}

		default:
			return fmt.Errorf("invalid 'additionalItems' definition: %#v", y)

		}
	}

	return nil
}

func (v *itemsValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return
	}

	if v.item != nil {
		for i, l := 0, len(y); i < l; i++ {
			err := ctx.ValidateValueWith(y[i], v.item)
			if err != nil {
				ctx.Report(&ErrInvalidItem{i, err})
			}
		}

		// no additionalItems

		return
	}

	if len(v.items) > 0 {
		var (
			i  = 0
			la = len(y)
			lb = len(v.items)
		)

		for ; i < la && i < lb; i++ {
			err := ctx.ValidateValueWith(y[i], v.items[i])
			if err != nil {
				ctx.Report(&ErrInvalidItem{i, err})
			}
		}

		// additionalItems
		if v.additionalItem == additionalItemsDenied {
			for ; i < la; i++ {
				ctx.Report(&ErrInvalidItem{i, fmt.Errorf("additional item is not allowed")})
			}
		} else if v.additionalItem != nil {
			for ; i < la; i++ {
				err := ctx.ValidateValueWith(y[i], v.additionalItem)
				if err != nil {
					ctx.Report(&ErrInvalidItem{i, err})
				}
			}
		}

		return
	}
}
