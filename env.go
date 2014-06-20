package jsonschema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type Env struct {
	// schema definition for schemas that can be interpreted by this environment
	SchemaSchema *Schema

	schemas  map[string]*Schema
	keywords map[string]keyword
}

type keyword struct {
	keyword   string
	priority  int
	prototype reflect.Type
}

func NewEnv() *Env {
	return &Env{
		schemas:  map[string]*Schema{},
		keywords: map[string]keyword{},
	}
}

func (e *Env) RegisterKeyword(key string, priority int, v Validator) {
	if _, found := e.keywords[key]; found {
		panic("keyword is already registered")
	}
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr {
		panic("Validator must be a pointer")
	}
	e.keywords[key] = keyword{key, priority, rt.Elem()}
}

func (e *Env) RegisterSchema(id string, data []byte) (*Schema, error) {
	schema, err := e.BuildSchema(id, data)
	if err != nil {
		return nil, err
	}

	e.schemas[schema.Id.String()] = schema
	return schema, nil
}

func (e *Env) BuildSchema(id string, data []byte) (*Schema, error) {
	var (
		obj map[string]interface{}
	)

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}

	if e.SchemaSchema != nil {
		err := e.SchemaSchema.Validate(obj)
		if err != nil {
			return nil, err
		}
	}

	builder := newBuilder(e)
	schema, err := builder.Build(id, obj)
	if err != nil {
		return nil, err
	}

	err = builder.resolve()
	if err != nil {
		return nil, err
	}

	if id != "" && normalizeRef(schema.Id.String()) != id {
		return nil, fmt.Errorf("schema id dit not match url (%q != %q)", id, schema.Id)
	}

	return schema, nil
}
