package jsonschema

import (
	"fmt"
)

var additionalItemsDenied = &Schema{}

type additionalItemsValidator struct {
	item *Schema
}

func (v *additionalItemsValidator) Setup(x interface{}, builder Builder) error {
	switch y := x.(type) {

	case map[string]interface{}:
		s, err := builder.Build("/additionalItems", y)
		if err != nil {
			return err
		}
		v.item = s
		return nil

	case bool:
		if !y {
			v.item = additionalItemsDenied
		}
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

	ctx.AdditionalItems = v.item
}
