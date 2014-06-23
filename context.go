package jsonschema

type Context struct {
	Type PrimitiveType

	value   interface{}
	errors  []error
	results map[string]error
}

func (c *Context) Report(err error) {
	c.errors = append(c.errors, err)
}

func (c *Context) ValidateWith(schema *Schema) error {
	if schema.RefSchema != nil {
		return c.ValidateWith(schema.RefSchema)
	}

	id := normalizeRef(schema.Id.String())
	if err, found := c.results[id]; found {
		return err
	}

	c.results[id] = nil

	var (
		ctx Context
		err error
	)

	ctx.value = c.value
	ctx.results = c.results

	for _, validator := range schema.Validators {
		validator.Validate(ctx.value, &ctx)
	}

	if len(ctx.errors) > 0 {
		err = &ErrInvalidInstance{schema, ctx.errors}
		c.results[id] = err
	}

	return err
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
