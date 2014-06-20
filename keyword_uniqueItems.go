package jsonschema

import (
	"fmt"
	"reflect"
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
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			if reflect.DeepEqual(y[i], y[j]) {
				ctx.Report(&ErrNotUnique{i, j, y[i]})
			}
		}
	}
}
