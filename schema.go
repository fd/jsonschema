package jsonschema

import (
	"bytes"
	"encoding/json"
	"net/url"
)

type Schema struct {
	Id         *url.URL
	Ref        *url.URL
	RefSchema  *Schema
	Validators []Validator
	Definition map[string]interface{}
	Subschemas map[string]*Schema
}

type Validator interface {
	Setup(b Builder) error
	Validate(interface{}, *Context)
}

type FormatValidator interface {
	IsValid(interface{}) bool
}

func (s *Schema) Validate(v interface{}) error {
	return newContext().ValidateValueWith(v, s)
}

func (s *Schema) ValidateData(d []byte) error {
	var (
		v interface{}
	)

	dec := json.NewDecoder(bytes.NewReader(d))
	dec.UseNumber()
	err := dec.Decode(&v)
	if err != nil {
		return err
	}

	return s.Validate(v)
}
