package jsonschema

import (
	"fmt"
)

type anyOfValidator struct {
	schemas []*Schema
}

func (v *anyOfValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("anyOf"); found {
		y, ok := x.([]interface{})
		if !ok || y == nil {
			return fmt.Errorf("invalid 'anyOf' definition: %#v", x)
		}

		schemas := make([]*Schema, len(y))
		for i, a := range y {
			b, ok := a.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid 'anyOf' definition: %#v", x)
			}

			schema, err := builder.Build(fmt.Sprintf("/anyOf/%d", i), b)
			if err != nil {
				return err
			}

			schemas[i] = schema
		}

		v.schemas = schemas
	}
	return nil
}

func (v *anyOfValidator) Validate(x interface{}, ctx *Context) {
	var (
		errors []error
	)

	for i, schema := range v.schemas {
		err := ctx.ValidateSelfWith(schema)
		if err == nil {
			return
		}

		if errors == nil {
			errors = make([]error, len(v.schemas))
		}

		errors[i] = err
	}

	ctx.Report(&ErrNotAnyOf{x, v.schemas, errors})
}
