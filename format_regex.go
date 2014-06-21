package jsonschema

import (
	"regexp"
)

type regexFormat struct{}

func (*regexFormat) IsValid(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return true
	}

	_, err := regexp.Compile(s)
	return err == nil
}
