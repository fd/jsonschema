package jsonschema

import (
	"fmt"
)

type dependenciesValidator struct {
	dependencies map[string]interface{}
}

func (v *dependenciesValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("dependencies"); found {
		y, ok := x.(map[string]interface{})
		if !ok || y == nil {
			return fmt.Errorf("invalid 'dependencies' definition: %#v", x)
		}

		dependencies := make(map[string]interface{}, len(y))
		for dependant, a := range y {
			switch b := a.(type) {
			case []interface{}:
				deps := make([]string, len(b))
				for i, d := range b {
					if e, ok := d.(string); ok {
						deps[i] = e
					} else {
						return fmt.Errorf("invalid 'dependencies' definition: %#v", x)
					}
				}
				dependencies[dependant] = deps

			case map[string]interface{}:
				schema, err := builder.Build("/dependencies/"+escapeJSONPointer(dependant), b)
				if err != nil {
					return err
				}
				dependencies[dependant] = schema

			default:
				return fmt.Errorf("invalid 'dependencies' definition: %#v", x)
			}
		}

		v.dependencies = dependencies
	}
	return nil
}

func (v *dependenciesValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return
	}

	for k, a := range v.dependencies {
		if _, found := y[k]; !found {
			continue
		}

		switch d := a.(type) {
		case []string:
			for _, dep := range d {
				if _, found := y[dep]; !found {
					ctx.Report(&ErrInvalidDependency{Property: k, Dependency: dep})
				}
			}

		case *Schema:
			err := ctx.ValidateValueWith(x, d)
			if err != nil {
				ctx.Report(&ErrInvalidDependency{Property: k, Schema: d, Err: err})
			}

		}
	}
}
