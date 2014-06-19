package jsonschema

import (
	"fmt"
	"math"
	"reflect"
)

type NotMultipleOfError struct {
	intFactor   int64
	floatFactor float64
	was         reflect.Value
}

func (e *NotMultipleOfError) Error() string {
	if e.intFactor > 0 {
		return fmt.Sprintf("expected %#v to be a multiple of %v", e.was.Interface(), e.intFactor)
	} else {
		return fmt.Sprintf("expected %#v to be a multiple of %v", e.was.Interface(), e.floatFactor)
	}
}

type multipleOfValidator struct {
	intFactor   int64
	floatFactor float64
}

func (v *multipleOfValidator) Setup(x interface{}, e *Env) error {
	xv := reflect.ValueOf(x)

	switch xv.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		i := xv.Int()
		if i <= 0 {
			return fmt.Errorf("multipleOf must be greater than 0 (was: %v)", i)
		}
		v.intFactor = int64(i)
		v.floatFactor = float64(i)

	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		i := xv.Uint()
		if i <= 0 {
			return fmt.Errorf("multipleOf must be greater than 0 (was: %v)", i)
		}
		if i > math.MaxInt64 {
			return fmt.Errorf("multipleOf is too large %v", i)
		}
		v.intFactor = int64(i)
		v.floatFactor = float64(i)

	case
		reflect.Float32,
		reflect.Float64:

		i := xv.Float()
		if i < math.SmallestNonzeroFloat64 {
			return fmt.Errorf("multipleOf must be greater than 0.0 (was: %v)", i)
		}
		v.floatFactor = i
		if math.Trunc(i) == i {
			if i > float64(math.MaxInt64) {
				return fmt.Errorf("multipleOf is too large %v", i)
			}
			v.intFactor = int64(i)
		}

	default:
		return fmt.Errorf("invalid multipleOf definition: %#v", x)

	}

	return nil
}

func (v *multipleOfValidator) Validate(x reflect.Value, ctx *Context) {
	if ctx.Type != IntegerType && ctx.Type != NumberType {
		return
	}

	for x.Kind() == reflect.Interface || x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	kind := x.Kind()

	ok := false
	if v.intFactor > 0 {
		if kind == reflect.Float32 ||
			kind == reflect.Float64 {
			f := x.Float()
			if math.Trunc(f) != f {
				ctx.Report(&NotMultipleOfError{v.intFactor, v.floatFactor, x})
				return
			}
			ok = int64(f)%v.intFactor == 0

		} else if kind == reflect.Int ||
			kind == reflect.Int8 ||
			kind == reflect.Int16 ||
			kind == reflect.Int32 ||
			kind == reflect.Int64 {
			ok = x.Int()%v.intFactor == 0

		} else if kind == reflect.Uint ||
			kind == reflect.Uint8 ||
			kind == reflect.Uint16 ||
			kind == reflect.Uint32 ||
			kind == reflect.Uint64 {
			ok = x.Uint()%uint64(v.intFactor) == 0
		}

	} else {
		if kind == reflect.Float32 ||
			kind == reflect.Float64 {
			ok = math.Remainder(x.Float(), v.floatFactor) < math.SmallestNonzeroFloat64

		} else if kind == reflect.Int ||
			kind == reflect.Int8 ||
			kind == reflect.Int16 ||
			kind == reflect.Int32 ||
			kind == reflect.Int64 {
			f := float64(x.Int())
			ok = math.Remainder(f, v.floatFactor) < math.SmallestNonzeroFloat64

		} else if kind == reflect.Uint ||
			kind == reflect.Uint8 ||
			kind == reflect.Uint16 ||
			kind == reflect.Uint32 ||
			kind == reflect.Uint64 {
			f := float64(x.Uint())
			ok = math.Remainder(f, v.floatFactor) < math.SmallestNonzeroFloat64

		}

	}

	if !ok {
		ctx.Report(&NotMultipleOfError{v.intFactor, v.floatFactor, x})
	}
}
