package jsonschema

import (
	"fmt"
	"reflect"
	"regexp"
)

type InvalidPatternError struct {
	pattern string
	was     reflect.Value
}

func (e *InvalidPatternError) Error() string {
	return fmt.Sprintf("expected %#v to be maych %q", e.was.Interface(), e.pattern)
}

type patternValidator struct {
	pattern string
	regexp  *regexp.Regexp
}

func (v *patternValidator) Setup(x interface{}, e *Env) error {
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

func (v *patternValidator) Validate(x reflect.Value, ctx *Context) {
	if !isString(x) {
		return
	}

	for x.Kind() == reflect.Interface || x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	ok := false

	if x.Kind() == reflect.String {
		ok = v.regexp.MatchString(x.String())
	} else if x.Kind() == reflect.Slice {
		ok = v.regexp.Match(x.Interface().([]byte))
	}

	if !ok {
		ctx.Report(&InvalidPatternError{v.pattern, x})
	}
}
