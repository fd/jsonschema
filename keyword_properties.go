package jsonschema

import (
	"fmt"
	"regexp"
)

type ErrInvalidProperty struct {
	Property string
	Err      error
}

func (e *ErrInvalidProperty) Error() string {
	return fmt.Sprintf("Invalid property %q: %s", e.Property, e.Err)
}

type propertiesValidator struct {
	properties           map[string]*Schema
	patterns             []*patternProperty
	additionalProperties *Schema
}

type patternProperty struct {
	pattern string
	regexp  *regexp.Regexp
	schema  *Schema
}

var additionalPropertiesDenied = &Schema{}

func (v *propertiesValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("properties"); found {
		defs, ok := x.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid 'properties' definition: %#v", x)
		}

		properties := make(map[string]*Schema, len(defs))
		for k, y := range defs {
			mdef, ok := y.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid 'properties' definition: %#v", x)
			}

			schema, err := builder.Build("/properties/"+escapeJSONPointer(k), mdef)
			if err != nil {
				return err
			}
			properties[k] = schema
		}

		v.properties = properties
	}

	if x, found := builder.GetKeyword("patternProperties"); found {
		defs, ok := x.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid 'patternProperties' definition: %#v", x)
		}

		patterns := make([]*patternProperty, 0, len(defs))
		for k, y := range defs {
			mdef, ok := y.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid 'patternProperties' definition: %#v", x)
			}

			reg, err := regexp.Compile(k)
			if err != nil {
				return fmt.Errorf("invalid 'patternProperties' definition: %#v (%s)", x, err)
			}

			schema, err := builder.Build("/patternProperties/"+escapeJSONPointer(k), mdef)
			if err != nil {
				return err
			}

			patterns = append(patterns, &patternProperty{k, reg, schema})
		}

		v.patterns = patterns
	}

	if x, ok := builder.GetKeyword("additionalProperties"); ok {
		switch y := x.(type) {

		case map[string]interface{}:
			s, err := builder.Build("/additionalProperties", y)
			if err != nil {
				return err
			}
			v.additionalProperties = s

		case bool:
			if !y {
				v.additionalProperties = additionalPropertiesDenied
			}

		default:
			return fmt.Errorf("invalid 'additionalProperties' definition: %#v", y)

		}
	}

	return nil
}

func (v *propertiesValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return
	}

	var pending = make(map[string]interface{}, len(y))
	for k, m := range y {
		pending[k] = m
	}

	// properties
	for k, schema := range v.properties {
		if m, found := pending[k]; found {
			delete(pending, k)

			err := schema.Validate(m)
			if err != nil {
				ctx.Report(&ErrInvalidProperty{k, err})
			}
		}
	}

	// patternProperties
	for _, pattern := range v.patterns {
		for k, m := range pending {
			if pattern.regexp.MatchString(k) {
				delete(pending, k)

				err := pattern.schema.Validate(m)
				if err != nil {
					ctx.Report(&ErrInvalidProperty{k, err})
				}
			}
		}
	}

	// additionalProperties
	if v.additionalProperties == additionalPropertiesDenied {
		for k := range pending {
			ctx.Report(&ErrInvalidProperty{k, fmt.Errorf("additional property is not allowed")})
		}

	} else if v.additionalProperties != nil {
		for k, m := range pending {
			err := v.additionalProperties.Validate(m)
			if err != nil {
				ctx.Report(&ErrInvalidProperty{k, err})
			}
		}

	}

}
