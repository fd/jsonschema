package jsonschema

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrNotAllOf struct {
	Value   interface{}
	Schemas []*Schema
	Errors  []error
}

func (e *ErrNotAllOf) Error() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "value must be all of:")

	for i, schema := range e.Schemas {
		var (
			err    = e.Errors[i]
			errstr = "<nil>"
		)

		if err != nil {
			errstr = strings.Replace(err.Error(), "\n", "\n    ", -1)
		}

		fmt.Fprintf(&buf, "\n- schema: %v\n  error:\n    %v", schema, errstr)
	}

	return buf.String()
}

type allOfValidator struct {
	schemas []*Schema
}

func (v *allOfValidator) Setup(x interface{}, e *Env) error {
	y, ok := x.([]interface{})
	if !ok || y == nil {
		return fmt.Errorf("invalid 'allOf' definition: %#v", x)
	}

	schemas := make([]*Schema, len(y))
	for i, a := range y {
		b, ok := a.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid 'allOf' definition: %#v", x)
		}

		schema, err := e.BuildSchema(b)
		if err != nil {
			return err
		}

		schemas[i] = schema
	}

	v.schemas = schemas
	return nil
}

func (v *allOfValidator) Validate(x interface{}, ctx *Context) {
	var (
		errors = make([]error, len(v.schemas))
		failed = false
	)

	for i, schema := range v.schemas {
		err := schema.Validate(x)

		if err != nil {
			failed = true
		}

		errors[i] = err
	}

	if failed {
		ctx.Report(&ErrNotAllOf{x, v.schemas, errors})
	}
}
