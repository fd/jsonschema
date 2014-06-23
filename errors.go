package jsonschema

import (
	"bytes"
	"fmt"
	"strings"
)

// ErrNotAllOf is returned when a `allOf` keyword failed.
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

// ErrNotAnyOf is returned when a `anyOf` keyword failed.
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

// ErrInvalidDependency is returned when a `dependency` keyword failed.
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

// ErrInvalidEnum is returned when a `enum` keyword failed.
type ErrInvalidEnum struct {
	Expected []interface{}
	Value    interface{}
}

func (e *ErrInvalidEnum) Error() string {
	return fmt.Sprintf("%v must be in %v", e.Value, e.Expected)
}

// ErrInvalidFormat is returned when a `format` keyword failed.
type ErrInvalidFormat struct {
	Value  interface{}
	Format string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("%#v did not match format '%s'", e.Value, e.Format)
}

// ErrInvalidItem is returned when a `item` keyword failed.
type ErrInvalidItem struct {
	Index int
	Err   error
}

func (e *ErrInvalidItem) Error() string {
	return fmt.Sprintf("Invalid item at %v: %s", e.Index, e.Err)
}

// ErrTooLarge is returned when a `maximum` keyword failed.
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

// ErrTooLong is returned when a `maxLength`, `maxItems` or a `maxProperties` keyword failed.
type ErrTooLong struct {
	max int
	was interface{}
}

func (e *ErrTooLong) Error() string {
	return fmt.Sprintf("expected len(%#v) to be smaller than %v", e.was, e.max)
}

// ErrTooSmall is returned when a `minimum` keyword failed.
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

// ErrTooShort is returned when a `minLength`, `minItems` or a `minProperties` keyword failed.
type ErrTooShort struct {
	min int
	was interface{}
}

func (e *ErrTooShort) Error() string {
	return fmt.Sprintf("expected len(%#v) to be larger than %v", e.was, e.min)
}

// ErrNotMultipleOf is returned when a `multipleOf` keyword failed.
type ErrNotMultipleOf struct {
	factor float64
	was    interface{}
}

func (e *ErrNotMultipleOf) Error() string {
	return fmt.Sprintf("expected %#v to be a multiple of %v", e.was, e.factor)
}

// ErrNotNot is returned when a `not` keyword failed.
type ErrNotNot struct {
	Value  interface{}
	Schema *Schema
}

func (e *ErrNotNot) Error() string {
	return fmt.Sprintf("value must not be valid for: %v", e.Schema)
}

// ErrNotOneOf is returned when a `oneOf` keyword failed.
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

// ErrInvalidPattern is returned when a `pattern` keyword failed.
type ErrInvalidPattern struct {
	pattern string
	was     string
}

func (e *ErrInvalidPattern) Error() string {
	return fmt.Sprintf("expected %#v to be maych %q", e.was, e.pattern)
}

// ErrInvalidProperty is returned when a `property` keyword failed.
type ErrInvalidProperty struct {
	Property string
	Err      error
}

func (e *ErrInvalidProperty) Error() string {
	return fmt.Sprintf("Invalid property %q: %s", e.Property, e.Err)
}

// ErrRequiredProperty is returned when a `required` keyword failed.
type ErrRequiredProperty struct {
	expected string
}

func (e *ErrRequiredProperty) Error() string {
	return fmt.Sprintf("missing required property: %q", e.expected)
}

// ErrInvalidType is returned when a `type` keyword failed.
type ErrInvalidType struct {
	expected []PrimitiveType
	was      interface{}
}

func (e *ErrInvalidType) Error() string {
	return fmt.Sprintf("expected type to be in %#v but was %#v", e.expected, e.was)
}

// ErrNotUnique is returned when a `uniqueItems` keyword failed.
type ErrNotUnique struct {
	IndexA int
	IndexB int
	Value  interface{}
}

func (e *ErrNotUnique) Error() string {
	return fmt.Sprintf("value at %d (%v) is not unique (repeated at %d)", e.IndexA, e.Value, e.IndexB)
}

// ErrInvalidInstance is returned when the instance is invalid.
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
