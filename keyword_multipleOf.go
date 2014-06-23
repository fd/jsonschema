package jsonschema

import (
	"encoding/json"
	"fmt"
	"math"
)

type multipleOfValidator struct {
	factor float64
}

func (v *multipleOfValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("multipleOf"); found {
		y, ok := x.(json.Number)
		if !ok {
			return fmt.Errorf("invalid 'multipleOf' definition: %#v", x)
		}

		f, err := y.Float64()
		if err != nil {
			return fmt.Errorf("invalid 'multipleOf' definition: %#v (%s)", x, err)
		}

		if f < math.SmallestNonzeroFloat64 {
			return fmt.Errorf("invalid 'multipleOf' definition: %#v", x)
		}

		v.factor = f
	}
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

	rem := math.Abs(math.Remainder(f, v.factor))
	rem /= v.factor // normalize rem between 0.0 and 1.0
	ok = rem < 0.000000001

	if !ok {
		ctx.Report(&ErrNotMultipleOf{v.factor, x})
	}
}
