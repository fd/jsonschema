package jsonschema

import (
	"time"
)

type datetimeFormat struct{}

func (*datetimeFormat) IsValid(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return true
	}

	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}
