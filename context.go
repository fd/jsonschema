package jsonschema

type Context struct {
	Type             PrimitiveType
	ExclusiveMaximum bool
	ExclusiveMinimum bool
	NextItem         int
	errors           []error
}

func (c *Context) Report(err error) {
	c.errors = append(c.errors, err)
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
