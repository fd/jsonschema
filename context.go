package jsonschema

type Context struct {
	errors []error
}

func (c *Context) Report(err error) {
	c.errors = append(c.errors, err)
}
