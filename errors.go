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

type ErrNotAnyOf struct {
	Value   interface{}
	Schemas []*Schema
	Errors  []error
}

func (e *ErrNotAnyOf) Error() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "value must be any of:")

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

type ErrInvalidDependency struct {
	Property   string
	Dependency string
	Schema     *Schema
	Err        error
}

func (e *ErrInvalidDependency) Error() string {
	if e.Schema != nil {
		return fmt.Sprintf("Invalid property %q: faile to validate dependecy: %s", e.Property, e.Err)
	} else {
		return fmt.Sprintf("Invalid property %q: missing property: %q", e.Property, e.Dependency)
	}
}

type ErrInvalidEnum struct {
	Expected []interface{}
	Value    interface{}
}

func (e *ErrInvalidEnum) Error() string {
	return fmt.Sprintf("%v must be in %v", e.Value, e.Expected)
}

type ErrInvalidFormat struct {
	Value  interface{}
	Format string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("%#v did not match format '%s'", e.Value, e.Format)
}

type ErrInvalidItem struct {
	Index int
	Err   error
}

func (e *ErrInvalidItem) Error() string {
	return fmt.Sprintf("Invalid item at %v: %s", e.Index, e.Err)
}

type ErrTooLarge struct {
	max       float64
	exclusive bool
	was       interface{}
}

func (e *ErrTooLarge) Error() string {
	if e.exclusive {
		return fmt.Sprintf("expected %#v to be smaller than %v", e.was, e.max)
	} else {
		return fmt.Sprintf("expected %#v to be smaller than or equal to %v", e.was, e.max)
	}
}

type ErrTooLong struct {
	max int
	was interface{}
}

func (e *ErrTooLong) Error() string {
	return fmt.Sprintf("expected len(%#v) to be smaller than %v", e.was, e.max)
}

type ErrTooSmall struct {
	min       float64
	exclusive bool
	was       interface{}
}

func (e *ErrTooSmall) Error() string {
	if e.exclusive {
		return fmt.Sprintf("expected %#v to be larger than %v", e.was, e.min)
	} else {
		return fmt.Sprintf("expected %#v to be larger than or equal to %v", e.was, e.min)
	}
}

type ErrTooShort struct {
	min int
	was interface{}
}

func (e *ErrTooShort) Error() string {
	return fmt.Sprintf("expected len(%#v) to be larger than %v", e.was, e.min)
}

type ErrNotMultipleOf struct {
	factor float64
	was    interface{}
}

func (e *ErrNotMultipleOf) Error() string {
	return fmt.Sprintf("expected %#v to be a multiple of %v", e.was, e.factor)
}

type ErrNotNot struct {
	Value  interface{}
	Schema *Schema
}

func (e *ErrNotNot) Error() string {
	return fmt.Sprintf("value must not be valid for: %v", e.Schema)
}

type ErrNotOneOf struct {
	Value   interface{}
	Schemas []*Schema
	Errors  []error
}

func (e *ErrNotOneOf) Error() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "value must be one of:")

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

type ErrInvalidPattern struct {
	pattern string
	was     string
}

func (e *ErrInvalidPattern) Error() string {
	return fmt.Sprintf("expected %#v to be maych %q", e.was, e.pattern)
}

type ErrInvalidProperty struct {
	Property string
	Err      error
}

func (e *ErrInvalidProperty) Error() string {
	return fmt.Sprintf("Invalid property %q: %s", e.Property, e.Err)
}

type ErrRequiredProperty struct {
	expected string
}

func (e *ErrRequiredProperty) Error() string {
	return fmt.Sprintf("missing required property: %q", e.expected)
}

type ErrInvalidType struct {
	expected []PrimitiveType
	was      interface{}
}

func (e *ErrInvalidType) Error() string {
	return fmt.Sprintf("expected type to be in %#v but was %#v", e.expected, e.was)
}

type ErrNotUnique struct {
	IndexA int
	IndexB int
	Value  interface{}
}

func (e *ErrNotUnique) Error() string {
	return fmt.Sprintf("value at %d (%v) is not unique (repeated at %d)", e.IndexA, e.Value, e.IndexB)
}

type ErrInvalidInstance struct {
	Schema *Schema
	Errors []error
}

func (e *ErrInvalidInstance) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Schema errors (%s):", normalizeRef(e.Schema.Id.String()))
	for _, err := range e.Errors {
		s := strings.Replace(err.Error(), "\n", "\n  ", -1)
		fmt.Fprintf(&buf, "\n- %s", s)
	}
	return buf.String()
}
