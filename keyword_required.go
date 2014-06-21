package jsonschema

import (
	"fmt"
)

type ErrRequiredProperty struct {
	expected string
}

func (e *ErrRequiredProperty) Error() string {
	return fmt.Sprintf("missing required property: %q", e.expected)
}

type requiredValidator struct {
	required []string
}

func (v *requiredValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("required"); found {
		switch y := x.(type) {

		case []string:
			v.required = y

		case []interface{}:
			var z = make([]string, len(y))
			for i, a := range y {
				if b, ok := a.(string); ok {
					z[i] = b
				} else {
					return fmt.Errorf("invalid 'required' definition: %#v", x)
				}
			}
			v.required = z

		default:
			return fmt.Errorf("invalid 'required' definition: %#v", x)
		}
	}
	return nil
}

func (v *requiredValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(map[string]interface{})
	if !ok || y == nil {
		return
	}

	for _, k := range v.required {
		_, found := y[k]
		if !found {
			ctx.Report(&ErrRequiredProperty{k})
		}
	}
}
