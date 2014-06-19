package jsonschema

import (
	"fmt"
	"reflect"
)

type LargerThanMaximumError struct {
	max       float64
	exclusive bool
	was       reflect.Value
}

func (e *LargerThanMaximumError) Error() string {
	if e.exclusive {
		return fmt.Sprintf("expected %#v to be smaller than %v", e.was.Interface(), e.max)
	} else {
		return fmt.Sprintf("expected %#v to be smaller than or equal to %v", e.was.Interface(), e.max)
	}
}

type maximumValidator struct {
	max float64
}

func (v *maximumValidator) Setup(x interface{}, e *Env) error {
	xv := reflect.ValueOf(x)

	switch xv.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		i := xv.Int()
		v.max = float64(i)

	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		i := xv.Uint()
		v.max = float64(i)

	case
		reflect.Float32,
		reflect.Float64:

		v.max = xv.Float()

	default:
		return fmt.Errorf("invalid 'maximum' definition: %#v", x)

	}

	return nil
}

func (v *maximumValidator) Validate(x reflect.Value, ctx *Context) {
	if !isInteger(x) && !isFloat(x) {
		return
	}

	for x.Kind() == reflect.Interface || x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	var (
		f  float64
		ok bool
	)

	switch x.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		f = float64(x.Int())

	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		f = float64(x.Uint())

	case
		reflect.Float32,
		reflect.Float64:
		f = x.Float()

	}

	if ctx.ExclusiveMaximum {
		ok = f < v.max
	} else {
		ok = f <= v.max
	}
	fmt.Printf("excl=%v max=%v f=%v ok=%v\n", ctx.ExclusiveMaximum, v.max, f, ok)

	if !ok {
		ctx.Report(&LargerThanMaximumError{v.max, ctx.ExclusiveMaximum, x})
	}
}
