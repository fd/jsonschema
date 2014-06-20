package jsonschema

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type ErrTooLong struct {
	max int
	was interface{}
}

func (e *ErrTooLong) Error() string {
	return fmt.Sprintf("expected len(%#v) to be smaller than %v", e.was, e.max)
}

type maxLengthValidator struct {
	max int
}

func (v *maxLengthValidator) Setup(x interface{}, e *Env) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'maxLength' definition: %#v", x)
	}

	i, err := y.Int64()
	if err != nil {
		return fmt.Errorf("invalid 'maxLength' definition: %#v (%s)", x, err)
	}

	v.max = int(i)
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
