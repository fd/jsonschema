package jsonschema

import (
	"fmt"
)

type ErrNotUnique struct {
	IndexA int
	IndexB int
	Value  interface{}
}

func (e *ErrNotUnique) Error() string {
	return fmt.Sprintf("value at %d (%v) is not unique (repeated at %d)", e.IndexA, e.Value, e.IndexB)
}

type uniqueItemsValidator struct {
	unique bool
}

func (v *uniqueItemsValidator) Setup(x interface{}, e *Env) error {
	if y, ok := x.(bool); ok && y {
		v.unique = true
	}

	return nil
}

func (v *uniqueItemsValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return
	}

	l := len(y)
	skip := make(map[int]bool, l)
	for i := 0; i < l; i++ {
		if skip[i] {
			continue
		}
		for j := i + 1; j < l; j++ {
			if skip[j] {
				continue
			}

			a, b := y[i], y[j]

			equal, err := isEqual(a, b)
			if err != nil {
				skip[j] = true
				ctx.Report(err)
				continue
			}

			// other values
			if equal {
				skip[j] = true
				ctx.Report(&ErrNotUnique{i, j, a})
			}
		}
	}
}
