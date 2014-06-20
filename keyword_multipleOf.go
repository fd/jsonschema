package jsonschema

import (
	"encoding/json"
	"fmt"
	"math"
)

type ErrNotMultipleOf struct {
	factor float64
	was    interface{}
}

func (e *ErrNotMultipleOf) Error() string {
	return fmt.Sprintf("expected %#v to be a multiple of %v", e.was, e.factor)
}

type multipleOfValidator struct {
	factor float64
}

func (v *multipleOfValidator) Setup(x interface{}, e *Env) error {
	y, ok := x.(json.Number)
	if !ok {
		return fmt.Errorf("invalid 'multipleOf' definition: %#v", x)
	}

	f, err := y.Float64()
	if err != nil {
		return fmt.Errorf("invalid 'multipleOf' definition: %#v (%s)", x, err)
	}

	v.factor = f
	return nil
}

func (v *multipleOfValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(json.Number)
	if !ok {
		return
	}

	f, err := y.Float64()
	if err != nil {
		ctx.Report(err)
		return
	}

	rem := math.Mod(f, v.factor)
	ok = rem < math.SmallestNonzeroFloat64
	fmt.Printf("f=%v factor=%v rem=%v ok=%v\n", f, v.factor, rem, ok)

	if !ok {
		ctx.Report(&ErrNotMultipleOf{v.factor, x})
	}
}
