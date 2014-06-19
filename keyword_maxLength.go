package jsonschema

import (
	"fmt"
	"reflect"
)

type LargerThanMaxLengthError struct {
	max int
	was reflect.Value
}

func (e *LargerThanMaxLengthError) Error() string {
	return fmt.Sprintf("expected len(%#v) to be smaller than %v", e.was.Interface(), e.max)
}

type maxLengthValidator struct {
	max int
}

func (v *maxLengthValidator) Setup(x interface{}, e *Env) error {
	xv := reflect.ValueOf(x)

	switch xv.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		i := xv.Int()
		v.max = int(i)

	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		i := xv.Uint()
		v.max = int(i)

	case
		reflect.Float32,
		reflect.Float64:

		if isIntegerFloat(xv) {
			v.max = int(xv.Float())
		} else {
			return fmt.Errorf("invalid 'maxLength' definition: %#v", x)
		}

	default:
		return fmt.Errorf("invalid 'maxLength' definition: %#v", x)

	}

	return nil
}

func (v *maxLengthValidator) Validate(x reflect.Value, ctx *Context) {
	if !isString(x) {
		return
	}

	for x.Kind() == reflect.Interface || x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Len() > v.max {
		ctx.Report(&LargerThanMaxLengthError{v.max, x})
	}
}
