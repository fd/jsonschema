package jsonschema

import (
	"fmt"
)

type definitionsValidator struct {
}

func (v *definitionsValidator) Setup(x interface{}, e *Env) error {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return fmt.Errorf("invalid 'definitions' definition: %#v", x)
	}

	schemas := make(map[string]*Schema, len(y))
	for name, a := range y {
		b, ok := a.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid 'definitions' definition: %#v", x)
		}

		schema, err := e.BuildSchema(b)
		if err != nil {
			return err
		}

		schemas[name] = schema
	}

	return nil
}

func (v *definitionsValidator) Validate(x interface{}, ctx *Context) {
}
