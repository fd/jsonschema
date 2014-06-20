package jsonschema

import (
	"fmt"
)

type ErrNotNot struct {
	Value  interface{}
	Schema *Schema
}

func (e *ErrNotNot) Error() string {
	return fmt.Sprintf("value must not be valid for: %v", e.Schema)
}

type notValidator struct {
	schema *Schema
}

func (v *notValidator) Setup(x interface{}, builder Builder) error {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return fmt.Errorf("invalid 'not' definition: %#v", x)
	}

	schema, err := builder.Build("/not", y)
	if err != nil {
		return err
	}

	v.schema = schema
	return nil
}

func (v *notValidator) Validate(x interface{}, ctx *Context) {
	err := ctx.ValidateWith(v.schema)
	if err == nil {
		ctx.Report(&ErrNotNot{x, v.schema})
	}
}
