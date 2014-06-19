package jsonschema

import (
	"fmt"
	"reflect"
)

type SmallerThanMinimumError struct {
	min       float64
	exclusive bool
	was       reflect.Value
}

func (e *SmallerThanMinimumError) Error() string {
	if e.exclusive {
		return fmt.Sprintf("expected %#v to be larger than %v", e.was.Interface(), e.min)
	} else {
		return fmt.Sprintf("expected %#v to be larger than or equal to %v", e.was.Interface(), e.min)
	}
}

type minimumValidator struct {
	min float64
}

func (v *minimumValidator) Setup(x interface{}, e *Env) error {
	xv := reflect.ValueOf(x)

	switch xv.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		i := xv.Int()
		v.min = float64(i)

	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		i := xv.Uint()
		v.min = float64(i)

	case
		reflect.Float32,
		reflect.Float64:

		v.min = xv.Float()

	default:
		return fmt.Errorf("invalid 'minimum' definition: %#v", x)

	}

	return nil
}

func (v *minimumValidator) Validate(x reflect.Value, ctx *Context) {
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

	if ctx.ExclusiveMinimum {
		ok = f > v.min
	} else {
		ok = f >= v.min
	}

	if !ok {
		ctx.Report(&SmallerThanMinimumError{v.min, ctx.ExclusiveMinimum, x})
	}
}
