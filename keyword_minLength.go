package jsonschema

import (
	"fmt"
	"reflect"
)

type LargerThanMinLengthError struct {
	min int
	was reflect.Value
}

func (e *LargerThanMinLengthError) Error() string {
	return fmt.Sprintf("expected len(%#v) to be larger than %v", e.was.Interface(), e.min)
}

type minLengthValidator struct {
	min int
}

func (v *minLengthValidator) Setup(x interface{}, e *Env) error {
	xv := reflect.ValueOf(x)

	switch xv.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		i := xv.Int()
		v.min = int(i)

	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		i := xv.Uint()
		v.min = int(i)

	case
		reflect.Float32,
		reflect.Float64:

		if isIntegerFloat(xv) {
			v.min = int(xv.Float())
		} else {
			return fmt.Errorf("invalid 'minLength' definition: %#v", x)
		}

	default:
		return fmt.Errorf("invalid 'minLength' definition: %#v", x)

	}

	return nil
}

func (v *minLengthValidator) Validate(x reflect.Value, ctx *Context) {
	if !isString(x) {
		return
	}

	for x.Kind() == reflect.Interface || x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Len() < v.min {
		ctx.Report(&LargerThanMinLengthError{v.min, x})
	}
}
