package jsonschema

import (
	"fmt"
)

type allOfValidator struct {
	schemas []*Schema
}

func (v *allOfValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("allOf"); found {
		y, ok := x.([]interface{})
		if !ok || y == nil {
			return fmt.Errorf("invalid 'allOf' definition: %#v", x)
		}

		schemas := make([]*Schema, len(y))
		for i, a := range y {
			b, ok := a.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid 'allOf' definition: %#v", x)
			}

			schema, err := builder.Build(fmt.Sprintf("/allOf/%d", i), b)
			if err != nil {
				return err
			}

			schemas[i] = schema
		}

		v.schemas = schemas
	}
	return nil
}

func (v *allOfValidator) Validate(x interface{}, ctx *Context) {
	var (
		errors = make([]error, len(v.schemas))
		failed = false
	)

	for i, schema := range v.schemas {
		err := ctx.ValidateWith(schema)

		if err != nil {
			failed = true
		}

		errors[i] = err
	}

	if failed {
		ctx.Report(&ErrNotAllOf{x, v.schemas, errors})
	}
}
