package jsonschema

import (
	"fmt"
)

type oneOfValidator struct {
	schemas []*Schema
}

func (v *oneOfValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("oneOf"); found {
		y, ok := x.([]interface{})
		if !ok || y == nil {
			return fmt.Errorf("invalid 'oneOf' definition: %#v", x)
		}

		schemas := make([]*Schema, len(y))
		for i, a := range y {
			b, ok := a.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid 'oneOf' definition: %#v", x)
			}

			schema, err := builder.Build(fmt.Sprintf("/oneOf/%d", i), b)
			if err != nil {
				return err
			}

			schemas[i] = schema
		}

		v.schemas = schemas
	}
	return nil
}

func (v *oneOfValidator) Validate(x interface{}, ctx *Context) {
	var (
		errors = make([]error, len(v.schemas))
		passed int
	)

	for i, schema := range v.schemas {
		err := ctx.ValidateWith(schema)

		if err == nil {
			passed++
		}

		errors[i] = err
	}

	if passed != 1 {
		ctx.Report(&ErrNotOneOf{x, v.schemas, errors})
	}
}
