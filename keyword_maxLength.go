package jsonschema

import (
	"fmt"
	"unicode/utf8"
)

type maxLengthValidator struct {
	max int
}

func (v *maxLengthValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("maxLength"); found {
		i, ok := x.(int64)
		if !ok {
			return fmt.Errorf("invalid 'maxLength' definition: %#v", x)
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
