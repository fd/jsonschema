package jsonschema

import (
	"fmt"
	"regexp"
)

type patternProperty struct {
	pattern string
	regexp  *regexp.Regexp
	schema  *Schema
}

type patternPropertiesValidator struct {
	patterns []*patternProperty
}

func (v *patternPropertiesValidator) Setup(builder Builder) error {
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
	return nil
}

func (v *patternPropertiesValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return
	}

	ctx.PatternProperties = v.patterns
}
