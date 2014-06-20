package jsonschema

import (
	"fmt"
)

type ErrInvalidDependency struct {
	Property   string
	Dependency string
	Schema     *Schema
	Err        error
}

func (e *ErrInvalidDependency) Error() string {
	if e.Schema != nil {
		return fmt.Sprintf("Invalid property %q: faile to validate dependecy: %s", e.Property, e.Err)
	} else {
		return fmt.Sprintf("Invalid property %q: missing property: %q", e.Property, e.Dependency)
	}
}

type dependenciesValidator struct {
	dependencies map[string]interface{}
}

func (v *dependenciesValidator) Setup(x interface{}, builder Builder) error {
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
			err := d.Validate(x)
			if err != nil {
				ctx.Report(&ErrInvalidDependency{Property: k, Schema: d, Err: err})
			}

		}
	}
}
