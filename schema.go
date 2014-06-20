package jsonschema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Schema struct {
	Id         *url.URL
	Ref        *url.URL
	RefSchema  *Schema
	Validators []Validator
	Definition map[string]interface{}
}

type Validator interface {
	Setup(x interface{}, b Builder) error
	Validate(interface{}, *Context)
}

func (s *Schema) Validate(v interface{}) error {
	var (
		ctx Context
	)

	ctx.value = v
	ctx.results = map[string]error{}
	return ctx.ValidateWith(s)
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

type InvalidDocumentError struct {
	Schema *Schema
	Errors []error
}

func (e *InvalidDocumentError) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "Schema errors:")
	for _, err := range e.Errors {
		s := strings.Replace(err.Error(), "\n", "\n  ", -1)
		fmt.Fprintf(&buf, "\n- %s", s)
	}
	return buf.String()
}
