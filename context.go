package jsonschema

import (
	"fmt"
)

type Context struct {
	stack []contextStackFrame
}

type contextStackFrame struct {
	valueId int
	value   interface{}
	errors  []error
	schema  *Schema
}

func newContext() *Context {
	return &Context{
		stack: make([]contextStackFrame, 0, 8),
	}
}

func (c *Context) Report(err error) {
	l := len(c.stack)
	frame := &c.stack[l-1]
	frame.errors = append(frame.errors, err)
}

func (c *Context) UpdateValue(x interface{}) {
	l := len(c.stack)
	frame := &c.stack[l-1]
	frame.value = x
}

func (c *Context) CurrentSchema() *Schema {
	l := len(c.stack)
	frame := &c.stack[l-1]
	return frame.schema
}

func (c *Context) ValidateValueWith(x interface{}, schema *Schema) (interface{}, error) {
	l := len(c.stack)

	if l == cap(c.stack) {
		tmp := make([]contextStackFrame, l, l*2)
		copy(tmp, c.stack)
		c.stack = tmp
	}

	if schema.RefSchema != nil {
		return c.ValidateValueWith(x, schema.RefSchema)
	}

	var (
		err         error
		parentFrame *contextStackFrame
		valueId     = 0
	)

	if l > 0 {
		parentFrame = &c.stack[l-1]
		valueId = parentFrame.valueId + 1
	}

	// push stack frame
	c.stack = append(c.stack, contextStackFrame{
		valueId: valueId,
		value:   x,
		schema:  schema,
	})

	for _, validator := range schema.Validators {
		validator.Validate(c.stack[l].value, c)
	}

	frame := &c.stack[l]
	if len(frame.errors) > 0 {
		err = &ErrInvalidInstance{schema, frame.errors}
	}

	// pop stack frame
	c.stack = c.stack[:len(c.stack)-1]
	return frame.value, err
}

func (c *Context) ValidateSelfWith(schema *Schema) (interface{}, error) {
	l := len(c.stack)

	if l == cap(c.stack) {
		tmp := make([]contextStackFrame, l, l*2)
		copy(tmp, c.stack)
		c.stack = tmp
	}

	if l == 0 {
		return nil, fmt.Errorf("ValidateWith() cannot be a root frame")
	}

	if schema.RefSchema != nil {
		return c.ValidateSelfWith(schema.RefSchema)
	}

	var (
		err         error
		parentFrame = &c.stack[l-1]
	)

	for i := l - 1; i >= 0; i-- {
		frame := &c.stack[i]
		if frame.valueId != parentFrame.valueId {
			break
		}
		if schema == frame.schema {
			return nil, fmt.Errorf("schema validation loops are invalid")
		}
	}

	// push stack frame
	c.stack = append(c.stack, contextStackFrame{
		valueId: parentFrame.valueId,
		value:   parentFrame.value,
		schema:  schema,
	})

	for _, validator := range schema.Validators {
		validator.Validate(c.stack[l].value, c)
	}

	frame := &c.stack[l]
	if len(frame.errors) > 0 {
		err = &ErrInvalidInstance{schema, frame.errors}
	}

	// pop stack frame
	c.stack = c.stack[:len(c.stack)-1]
	return frame.value, err
}

type PrimitiveType string

const (
	ArrayType   = PrimitiveType("array")
	BooleanType = PrimitiveType("boolean")
	IntegerType = PrimitiveType("integer")
	NullType    = PrimitiveType("null")
	NumberType  = PrimitiveType("number")
	ObjectType  = PrimitiveType("object")
	StringType  = PrimitiveType("string")
)

func (p PrimitiveType) Valid() bool {
	return p == ArrayType ||
		p == BooleanType ||
		p == IntegerType ||
		p == NullType ||
		p == NumberType ||
		p == ObjectType ||
		p == StringType
}
