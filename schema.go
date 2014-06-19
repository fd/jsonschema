package jsonschema

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type Schema struct {
	Validators []Validator
	Definition map[string]interface{}
}

type Validator interface {
	Setup(x interface{}, e *Env) error
	Validate(reflect.Value, *Context)
}

func (s *Schema) Validate(v interface{}) error {
	return s.ValidateValue(reflect.ValueOf(v))
}

func (s *Schema) ValidateValue(v reflect.Value) error {
	var ctx Context

	for _, validator := range s.Validators {
		validator.Validate(v, &ctx)
	}

	if len(ctx.errors) > 0 {
		return &InvalidDocumentError{s, ctx.errors}
	}

	return nil
}

type InvalidDocumentError struct {
	Schema *Schema
	Errors []error
}

func (e *InvalidDocumentError) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "Schema errors:")
	for _, err := range e.Errors {
		s := strings.Replace(err.Error(), "\n", "\n  ", -1)
		fmt.Fprintf(&buf, "\n- %s", s)
	}
	return buf.String()
}
