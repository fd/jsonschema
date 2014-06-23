package jsonschema

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type minLengthValidator struct {
	min int
}

func (v *minLengthValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("minLength"); found {
		y, ok := x.(json.Number)
		if !ok {
			return fmt.Errorf("invalid 'minLength' definition: %#v", x)
		}

		i, err := y.Int64()
		if err != nil {
			return fmt.Errorf("invalid 'minLength' definition: %#v (%s)", x, err)
		}

		v.min = int(i)
	}
	return nil
}

func (v *minLengthValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(string)
	if !ok {
		return
	}

	l := utf8.RuneCountInString(y)

	if l < v.min {
		ctx.Report(&ErrTooShort{v.min, x})
	}
}
