package jsonschema

import (
	"fmt"
	"net/url"
)

type uriFormat struct{}

func (*uriFormat) IsValid(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return true
	}

	u, err := url.Parse(s)
	fmt.Printf("url=%#v\n", u)
	return err == nil
}
