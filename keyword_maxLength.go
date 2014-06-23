package jsonschema

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type maxLengthValidator struct {
	max int
}

func (v *maxLengthValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("maxLength"); found {
		y, ok := x.(json.Number)
		if !ok {
			return fmt.Errorf("invalid 'maxLength' definition: %#v", x)
		}

		i, err := y.Int64()
		if err != nil {
			return fmt.Errorf("invalid 'maxLength' definition: %#v (%s)", x, err)
		}

		v.max = int(i)
	}
	return nil
}

func (v *maxLengthValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(string)
	if !ok {
		return
	}

	l := utf8.RuneCountInString(y)

	if l > v.max {
		ctx.Report(&ErrTooLong{v.max, x})
	}
}
