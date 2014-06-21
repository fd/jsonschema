package jsonschema

import (
	"fmt"
)

type ErrInvalidFormat struct {
	Value  interface{}
	Format string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("%#v did not match format '%s'", e.Value, e.Format)
}

type formatValidator struct {
	name   string
	format FormatValidator
}

func (v *formatValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("format"); found {
		y, ok := x.(string)
		if !ok {
			return fmt.Errorf("invalid 'format' definition: %#v", x)
		}

		format := builder.GetFormatValidator(y)
		if format == nil {
			return fmt.Errorf("invalid 'format' definition: %#v (unknown format)", x)
		}

		v.name = y
		v.format = format
	}
	return nil
}

func (v *formatValidator) Validate(x interface{}, ctx *Context) {
	if !v.format.IsValid(x) {
		ctx.Report(&ErrInvalidFormat{x, v.name})
	}
}
