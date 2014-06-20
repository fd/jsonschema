package jsonschema

import (
	"fmt"
)

type ErrInvalidProperty struct {
	Property string
	Err      error
}

func (e *ErrInvalidProperty) Error() string {
	return fmt.Sprintf("Invalid property %q: %s", e.Property, e.Err)
}

type propertiesValidator struct {
	members map[string]*Schema
}

func (v *propertiesValidator) Setup(x interface{}, builder Builder) error {
	defs, ok := x.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid 'properties' definition: %#v", x)
	}

	members := make(map[string]*Schema, len(defs))
	for k, y := range defs {
		mdef, ok := y.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid 'properties' definition: %#v", x)
		}

		schema, err := builder.Build("/properties/"+escapeJSONPointer(k), mdef)
		if err != nil {
			return err
		}
		members[k] = schema
	}

	v.members = members
	return nil
}

func (v *propertiesValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return
	}

	for k, m := range y {
		schema, ok := v.members[k]
		if !ok {
			continue
		}

		err := schema.Validate(m)
		if err != nil {
			ctx.Report(&ErrInvalidProperty{k, err})
		}
	}
}
