package jsonschema

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type ErrTooShort struct {
	min int
	was interface{}
}

func (e *ErrTooShort) Error() string {
	return fmt.Sprintf("expected len(%#v) to be larger than %v", e.was, e.min)
}

type minLengthValidator struct {
	min int
}

func (v *minLengthValidator) Setup(x interface{}, e *Env) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'minLength' definition: %#v", x)
	}

	i, err := y.Int64()
	if err != nil {
		return fmt.Errorf("invalid 'minLength' definition: %#v (%s)", x, err)
	}

	v.min = int(i)
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
