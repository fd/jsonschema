package jsonschema

import (
	"fmt"
	"regexp"
)

type ErrInvalidPattern struct {
	pattern string
	was     string
}

func (e *ErrInvalidPattern) Error() string {
	return fmt.Sprintf("expected %#v to be maych %q", e.was, e.pattern)
}

type patternValidator struct {
	pattern string
	regexp  *regexp.Regexp
}

func (v *patternValidator) Setup(builder Builder) error {
	if x, found := builder.GetKeyword("pattern"); found {
		if y, ok := x.(string); ok {
			r, err := regexp.Compile(y)
			if err != nil {
				return fmt.Errorf("invalid 'pattern' definition: %#v (error: %s)", x, err)
			}
			v.pattern = y
			v.regexp = r
			return nil
		}

		return fmt.Errorf("invalid 'pattern' definition: %#v", x)
	}

	return nil
}

func (v *patternValidator) Validate(x interface{}, ctx *Context) {
	y, ok := x.(string)
	if !ok {
		return
	}

	if !v.regexp.MatchString(y) {
		ctx.Report(&ErrInvalidPattern{v.pattern, y})
	}
}
