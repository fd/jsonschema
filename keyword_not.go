package jsonschema

import (
	"fmt"
)

type notValidator struct {
	schema *Schema
}

func (v *notValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("not"); found {
		y, ok := x.(map[string]interface{})
		if !ok || y == nil {
			return fmt.Errorf("invalid 'not' definition: %#v", x)
		}

		schema, err := builder.Build("/not", y)
		if err != nil {
			return err
		}

		v.schema = schema
	}
	return nil
}

func (v *notValidator) Validate(x interface{}, ctx *Context) {
	_, err := ctx.ValidateSelfWith(v.schema)
	if err == nil {
		ctx.Report(&ErrNotNot{x, v.schema})
	}
}
